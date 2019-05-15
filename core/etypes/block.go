package etypes

import (
	"errors"
	"github.com/eager7/echain/core/pb"
	"github.com/golang/protobuf/proto"
)

type Block struct {
	Header
	CountTxs     uint64         `json:"count_txs"`
	Transactions []*Transaction `json:"transactions"`
}

func (b *Block) Serialize() ([]byte, error) {
	pBlock := &pb.Block{
		Header: &pb.Header{
			Version:     b.Version,
			ChainID:     b.ChainID.Bytes(),
			Timestamp:   b.TimeStamp,
			Height:      b.Height,
			PrevHash:    b.PrevHash.Bytes(),
			MerkleHash:  b.MerkleHash.Bytes(),
			StateHash:   b.StateHash.Bytes(),
			ReceiptHash: b.ReceiptHash.Bytes(),
			Bloom:       b.Bloom.Bytes(),
			Hash:        b.Hash.Bytes(),
			Sign: &pb.Signature{
				PubKey:  b.Signatures.PubKey,
				SigData: b.Signatures.SigData,
			},
		},
		TxsCount:     b.CountTxs,
		Transactions: []*pb.Transaction{},
	}
	for _, tx := range b.Transactions {
		pBlock.Transactions = append(pBlock.Transactions, &pb.Transaction{
			Version:   tx.Version,
			ChainID:   tx.ChainID.Bytes(),
			From:      tx.From.Bytes(),
			To:        tx.To.Bytes(),
			Payload:   tx.Payload,
			Timestamp: tx.TimeStamp,
			Sign: &pb.Signature{
				PubKey:  tx.Signatures.PubKey,
				SigData: tx.Signatures.SigData,
			},
			Hash: tx.Hash.Bytes(),
		})
	}
	d, err := proto.Marshal(pBlock)
	if err != nil {
		return nil, errors.New("marshal err:" + err.Error())
	}
	return d, nil
}

func (b *Block) Deserialize(data []byte) error {

	return nil
}

func (b *Block) Instance() interface{} {
	return b
}

func (b *Block) Identify() pb.Id {
	return pb.Id_BlockType
}
