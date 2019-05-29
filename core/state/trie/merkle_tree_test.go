package trie_test

import (
	"fmt"
	"github.com/eager7/echain/core/etypes"
	"github.com/eager7/echain/core/state/trie"
	"testing"
)

func TestMerkleTree(t *testing.T) {
	fmt.Println("test merkle tree")

	var hashes []etypes.Hash
	h := etypes.HashSetHex("E11E799DE47FE0BFF8776956C401FE4A9532F9F406A092D3F5606CF6A8E18AEE")
	hashes = append(hashes, h)
	merkleRoot, _ := trie.GetMerkleRoot(hashes)
	fmt.Println(merkleRoot.Hex())
}
