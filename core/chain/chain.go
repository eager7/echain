package chain

import (
	"encoding/binary"
	"errors"
	"github.com/eager7/echain/core/etypes"
	mpt "github.com/eager7/echain/core/state"
	"github.com/eager7/echain/core/state/store"
)

type Chain struct {
	BlockStore  store.Storage
	HeightStore store.Storage
	State       *mpt.Mpt
}

func Initialize(dir string) (*Chain, error) {
	bs, err := store.NewBlockStore(dir + "/block")
	if err != nil {
		return nil, errors.New("new block store err:" + err.Error())
	}
	hs, err := store.NewBlockStore(dir + "/height")
	if err != nil {
		return nil, errors.New("new height store err:" + err.Error())
	}
	state, err := mpt.NewMptTree(dir+"/state", etypes.Hash{})
	if err != nil {
		return nil, errors.New("new mpt err:" + err.Error())
	}
	return &Chain{
		BlockStore:  bs,
		HeightStore: hs,
		State:       state,
	}, nil
}

func (c *Chain) StoreBlock(block *etypes.Block) error {
	data, err := block.Serialize()
	if err != nil {
		return errors.New("store block err:" + err.Error())
	}
	if err := c.BlockStore.Put(block.Hash.Bytes(), data); err != nil {
		return errors.New("store block db err:" + err.Error())
	}
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, block.Height)
	b := buf[:n]
	if err := c.HeightStore.Put(b, block.Hash.Bytes()); err != nil {
		return errors.New("store height db err:" + err.Error())
	}
	return nil
}

func (c *Chain) QueryBlock(height uint64) (*etypes.Block, error) {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, height)
	key := buf[:n]
	hd, err := c.HeightStore.Get(key)
	if err != nil {
		return nil, errors.New("query block err[height get]:" + err.Error())
	}
	hb, err := c.BlockStore.Get(hd)
	if err != nil {
		return nil, errors.New("query block err[block get]:" + err.Error())
	}
	var block etypes.Block
	if err := block.Deserialize(hb); err != nil {
		return nil, errors.New("query block err[deserialize]:" + err.Error())
	}
	return &block, nil
}
