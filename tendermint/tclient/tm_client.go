package tclient

import (
	"errors"
	"fmt"
	"github.com/eager7/echain/core/etypes"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

type TmClient struct {
	rpc client.Client
}

func Initialize(sock string) *TmClient {
	cli := client.NewHTTP(sock, "/websocket")
	return &TmClient{rpc: cli}
}

func (t *TmClient) SendTx(tx *etypes.Transaction) (*etypes.Hash, error) {
	data, err := tx.Serialize()
	if err != nil {
		return nil, err
	}
	resp, err := t.rpc.BroadcastTxCommit(types.Tx(data))
	if err != nil {
		return nil, err
	}
	if resp.CheckTx.Code != 0 {
		return nil, errors.New("check tx error:" + resp.CheckTx.String())
	}
	if resp.DeliverTx.Code != 0 {
		return nil, errors.New("deliver tx error:" + resp.DeliverTx.String())
	}
	h := etypes.HashSetBytes(resp.Hash.Bytes())
	fmt.Println("tx hash:", h.Hex(), "echain hash:", tx.Hash.Hex())
	return &tx.Hash, nil
}

func (t *TmClient) QueryBlock(height int64) (*types.Block, []*etypes.Transaction, error) {
	resp, err := t.rpc.Block(&height)
	if err != nil {
		return nil, nil, errors.New("query block err[rpc block]:" + err.Error())
	}
	var tx []*etypes.Transaction
	for _, t := range resp.Block.Txs {
		var et etypes.Transaction
		if err := et.Deserialize(t); err != nil {
			return nil, nil, errors.New("query block err[tx deserialize]:" + err.Error())
		}
		tx = append(tx, &et)
	}
	return resp.Block, tx, nil
}
