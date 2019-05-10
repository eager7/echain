package types

import (
	"errors"
	"github.com/eager7/echain/core/pb"
	"github.com/golang/protobuf/proto"
	"math/big"
)

type Transaction struct {
	Version    uint32      `json:"version"`
	ChainID    Hash        `json:"chain_id"`
	From       Address     `json:"from"`
	To         Address     `json:"to"`
	Amount     *big.Int    `json:"amount"`
	TimeStamp  int64       `json:"timeStamp"`
	Payload    []byte      `json:"payload"`
	Signatures Signature `json:"signatures"`
	Hash       Hash        `json:"hash"`
}

func (t *Transaction) Serialize() ([]byte, error) {
	d, err := proto.Marshal(&pb.Transaction{
		Version:              t.Version,
		ChainID:              t.ChainID.Bytes(),
		From:                 t.From.Bytes(),
		Addr:                 t.To.Bytes(),
		Payload:              t.Payload,
		Timestamp:            t.TimeStamp,
		Sign:                 &pb.Signature{
			PubKey:               t.Signatures.PubKey,
			SigData:              t.Signatures.SigData,
		},
		Hash:                 t.Hash.Bytes(),
	})
	if err != nil {
		return nil, errors.New("marshal err:" + err.Error())
	}
	return d, nil
}

func (t *Transaction) Deserialize(data []byte) error {

	return nil
}

func (t *Transaction) Instance() interface{} {
	return t
}

func (t *Transaction) Identify() pb.Id {
	return pb.Id_TransactionType
}
