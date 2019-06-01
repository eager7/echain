package trie

import (
	"crypto/sha256"
	"github.com/eager7/echain/core/etypes"
)

type MerkleTree struct {
	Depth uint
	Root  *MerkleNode
}

/**
** Leaf Node
 */
type MerkleNode struct {
	hash  etypes.Hash
	left  *MerkleNode
	right *MerkleNode
}

func NewMerkleTree(hashes []etypes.Hash) *MerkleTree {
	if len(hashes) == 0 {
		return nil
	}
	var nodes []*MerkleNode
	for _, h := range hashes {
		nodes = append(nodes, &MerkleNode{h, nil, nil})
	}
	var height uint = 1
	for len(nodes) > 1 {
		nodes = buildTree(nodes)
		height += 1
	}
	return &MerkleTree{
		Depth: height,
		Root:  nodes[0],
	}
}

/**
** Create Merkle
 */
func buildTree(nodes []*MerkleNode) []*MerkleNode {
	var rootNode []*MerkleNode
	for i := 0; i < len(nodes)/2; i++ {
		var data []etypes.Hash
		data = append(data, nodes[i*2].hash)
		data = append(data, nodes[i*2+1].hash)
		hash := merkleHash(data)
		parentNode := &MerkleNode{
			hash:  hash,
			left:  nodes[i*2],
			right: nodes[i*2+1],
		}
		rootNode = append(rootNode, parentNode)
	}
	if len(nodes)%2 == 1 {
		var data []etypes.Hash
		data = append(data, nodes[len(nodes)-1].hash)
		data = append(data, nodes[len(nodes)-1].hash)
		hash := merkleHash(data)
		parentNode := &MerkleNode{
			hash:  hash,
			left:  nodes[len(nodes)-1],
			right: nodes[len(nodes)-1],
		}
		rootNode = append(rootNode, parentNode)
	}
	return rootNode
}

/**
** Compute parent hash
 */
func merkleHash(hashes []etypes.Hash) etypes.Hash {
	var joinHash []byte
	for _, h := range hashes {
		joinHash = append(joinHash, h.Bytes()...)
	}
	temp := sha256.Sum256(joinHash)
	f := sha256.Sum256(temp[:])
	return etypes.Hash(f)
}

/**
** compute merkle hash of hash list
 */
func GetMerkleRoot(hashes []etypes.Hash) (etypes.Hash, error) {
	if len(hashes) == 0 {
		return etypes.Hash{}, nil
	}
	if len(hashes) == 1 {
		hashes = append(hashes, hashes[0])
	}

	return NewMerkleTree(hashes).Root.hash, nil
}
