package trie

import (
	"fmt"
	"github.com/eager7/echain/core/etypes"
	"github.com/eager7/echain/core/state/store"
	"sync"
	"time"
)

const secureKeyLength = 11 + 32

var secureKeyPrefix = []byte("secure-key-")

type DatabaseReader interface {
	Get(key []byte) (value []byte, err error)
	Has(key []byte) (bool, error)
}

type Database struct {
	diskDB    store.Database              // Persistent storage for matured trie nodes
	nodes     map[etypes.Hash]*cachedNode // Data and references relationships of a node
	preImages map[etypes.Hash][]byte      // PreImages of nodes from the secure trie
	secKeyBuf [secureKeyLength]byte       // Ephemeral buffer for calculating preImage keys

	gcTime  time.Duration      // Time spent on garbage collection since last commit
	gcNodes uint64             // Nodes garbage collected since last commit
	gcSize  float64 // Data storage garbage collected since last commit

	nodesSize     float64 // Storage size of the nodes cache
	preImagesSize float64 // Storage size of the preImages cache

	lock sync.RWMutex
}

type cachedNode struct {
	blob     []byte              // Cached data block of the trie node
	parents  int                 // Number of live nodes referencing this one
	children map[etypes.Hash]int // Children referenced by this nodes
}

func NewDatabase(diskDB store.Database) *Database {
	return &Database{
		diskDB: diskDB,
		nodes: map[etypes.Hash]*cachedNode{
			{}: {children: make(map[etypes.Hash]int)},
		},
		preImages: make(map[etypes.Hash][]byte),
	}
}

func (db *Database) DiskDB() DatabaseReader {
	return db.diskDB
}

func (db *Database) Insert(hash etypes.Hash, blob []byte) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.insert(hash, blob)
}

func (db *Database) insert(hash etypes.Hash, blob []byte) {
	if _, ok := db.nodes[hash]; ok {
		return
	}
	db.nodes[hash] = &cachedNode{
		blob:     blob,
		children: make(map[etypes.Hash]int),
	}
	db.nodesSize += float64(etypes.HashLen + len(blob))
}

func (db *Database) insertPreImage(hash etypes.Hash, preImage []byte) {
	if _, ok := db.preImages[hash]; ok {
		return
	}
	db.preImages[hash] = preImage[:]
	db.preImagesSize += float64(etypes.HashLen + len(preImage))
}

func (db *Database) Node(hash etypes.Hash) ([]byte, error) {
	// Retrieve the node from cache if available
	db.lock.RLock()
	node := db.nodes[hash]
	db.lock.RUnlock()

	if node != nil {
		return node.blob, nil
	}
	// Content unavailable in memory, attempt to retrieve from disk
	return db.diskDB.Get(hash[:])
}

func (db *Database) preImage(hash etypes.Hash) ([]byte, error) {
	// Retrieve the node from cache if available
	db.lock.RLock()
	preImage := db.preImages[hash]
	db.lock.RUnlock()

	if preImage != nil {
		return preImage, nil
	}
	// Content unavailable in memory, attempt to retrieve from disk
	return db.diskDB.Get(db.secureKey(hash[:]))
}

func (db *Database) secureKey(key []byte) []byte {
	buf := append(db.secKeyBuf[:0], secureKeyPrefix...)
	buf = append(buf, key...)
	return buf
}

func (db *Database) Nodes() []etypes.Hash {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var hashes = make([]etypes.Hash, 0, len(db.nodes))
	for hash := range db.nodes {
		if hash != (etypes.Hash{}) { // Special case for "root" references/nodes
			hashes = append(hashes, hash)
		}
	}
	return hashes
}

func (db *Database) Reference(child etypes.Hash, parent etypes.Hash) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	db.reference(child, parent)
}

func (db *Database) reference(child etypes.Hash, parent etypes.Hash) {
	// If the node does not exist, it's a node pulled from disk, skip
	node, ok := db.nodes[child]
	if !ok {
		return
	}
	// If the reference already exists, only duplicate for roots
	if _, ok = db.nodes[parent].children[child]; ok && parent != (etypes.Hash{}) {
		return
	}
	node.parents++
	db.nodes[parent].children[child]++
}

func (db *Database) Dereference(child etypes.Hash, parent etypes.Hash) {
	db.lock.Lock()
	defer db.lock.Unlock()

	nodes, storage, start := len(db.nodes), db.nodesSize, time.Now()
	db.dereference(child, parent)

	db.gcNodes += uint64(nodes - len(db.nodes))
	db.gcSize += storage - db.nodesSize
	db.gcTime += time.Since(start)

	fmt.Println("DeReferenced trie from memory database", "nodes", nodes-len(db.nodes), "size", storage-db.nodesSize, "time", time.Since(start),
		"gcNodes", db.gcNodes, "gcSize", db.gcSize, "gcTime", db.gcTime, "liveNodes", len(db.nodes), "liveSize", db.nodesSize)
}

// dereference is the private locked version of Dereference.
func (db *Database) dereference(child etypes.Hash, parent etypes.Hash) {
	// Dereference the parent-child
	node := db.nodes[parent]

	node.children[child]--
	if node.children[child] == 0 {
		delete(node.children, child)
	}
	// If the node does not exist, it's a previously committed node.
	node, ok := db.nodes[child]
	if !ok {
		return
	}
	// If there are no more references to the child, delete it and cascade
	node.parents--
	if node.parents == 0 {
		for hash := range node.children {
			db.dereference(hash, child)
		}
		delete(db.nodes, child)
		db.nodesSize -= float64(etypes.HashLen + len(node.blob))
	}
}

func (db *Database) Commit(node etypes.Hash, report bool) error {
	db.lock.RLock()

	start := time.Now()
	batch := db.diskDB.NewBatch()

	for hash, preImage := range db.preImages {
		if err := batch.Put(db.secureKey(hash[:]), preImage); err != nil {
			fmt.Println("Failed to commit preImage from trie database", "err", err)
			db.lock.RUnlock()
			return err
		}
		if batch.ValueSize() > store.IdealBatchSize {
			if err := batch.Write(); err != nil {
				return err
			}
			batch.Reset()
		}
	}
	nodes, storage := len(db.nodes), db.nodesSize+db.preImagesSize
	if err := db.commit(node, batch); err != nil {
		fmt.Println("Failed to commit trie from trie database", "err", err)
		db.lock.RUnlock()
		return err
	}
	if err := batch.Write(); err != nil {
		fmt.Println("Failed to write trie to disk", "err", err)
		db.lock.RUnlock()
		return err
	}
	db.lock.RUnlock()

	db.lock.Lock()
	defer db.lock.Unlock()

	db.preImages = make(map[etypes.Hash][]byte)
	db.preImagesSize = 0

	db.unCache(node)

	if !report {
		fmt.Println("Persisted trie from memory database", "nodes", nodes-len(db.nodes), "size", storage-db.nodesSize, "time", time.Since(start),
			"gcNodes", db.gcNodes, "gcSize", db.gcSize, "gcTime", db.gcTime, "liveNodes", len(db.nodes), "liveSize", db.nodesSize)
	}

	db.gcNodes, db.gcSize, db.gcTime = 0, 0, 0

	return nil
}

func (db *Database) commit(hash etypes.Hash, batch store.Batch) error {
	node, ok := db.nodes[hash]
	if !ok {
		return nil
	}
	for child := range node.children {
		if err := db.commit(child, batch); err != nil {
			return err
		}
	}
	if err := batch.Put(hash[:], node.blob); err != nil {
		return err
	}
	if batch.ValueSize() >= store.IdealBatchSize {
		if err := batch.Write(); err != nil {
			return err
		}
		batch.Reset()
	}
	return nil
}

func (db *Database) unCache(hash etypes.Hash) {
	node, ok := db.nodes[hash]
	if !ok {
		return
	}
	for child := range node.children {
		db.unCache(child)
	}
	delete(db.nodes, hash)
	db.nodesSize -= float64(etypes.HashLen + len(node.blob))
}

func (db *Database) Size() float64 {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return db.nodesSize + db.preImagesSize
}
