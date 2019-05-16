package etypes

import (
	"errors"
	"github.com/eager7/echain/core/pb"
	"github.com/golang/protobuf/proto"
	"math/big"
)

type Transaction struct {
	Version    uint32    `json:"version"`
	ChainID    Hash      `json:"chain_id"`
	From       Address   `json:"from"`
	To         Address   `json:"to"`
	Amount     *big.Int  `json:"amount"`
	TimeStamp  int64     `json:"timeStamp"`
	Payload    []byte    `json:"payload"`
	Signatures Signature `json:"signatures"`
	Hash       Hash      `json:"hash"`
}

func (t *Transaction) Serialize() ([]byte, error) {
	amount, err := t.Amount.GobEncode()
	if err != nil {
		return nil, errors.New("amount encode err:" + err.Error())
	}
	d, err := proto.Marshal(&pb.Transaction{
		Version:   t.Version,
		ChainID:   t.ChainID.Bytes(),
		From:      t.From.Bytes(),
		To:        t.To.Bytes(),
		Amount:    amount,
		Payload:   t.Payload,
		Timestamp: t.TimeStamp,
		Sign: &pb.Signature{
			PubKey:  t.Signatures.PubKey,
			SigData: t.Signatures.SigData,
		},
		Hash: t.Hash.Bytes(),
	})
	if err != nil {
		return nil, errors.New("marshal err:" + err.Error())
	}
	return d, nil
}

func (t *Transaction) Deserialize(data []byte) error {
	var px pb.Transaction
	if err := proto.Unmarshal(data, &px); err != nil {
		return errors.New("unmarshal tx err:" + err.Error())
	}
	amount := new(big.Int)
	if err := amount.GobDecode(px.Amount); err != nil {
		return errors.New("amount decode err:" + err.Error())
	}
	t.Version = px.Version
	t.ChainID = HashSetBytes(px.Hash)
	t.From = AddressSetBytes(px.From)
	t.To = AddressSetBytes(px.To)
	t.Amount = amount
	t.TimeStamp = px.Timestamp
	t.Payload = px.Payload
	t.Signatures = Signature{PubKey: px.Sign.PubKey, SigData: px.Sign.SigData}
	t.Hash = HashSetBytes(px.Hash)
	return nil
}

func (t *Transaction) Instance() interface{} {
	return t
}

func (t *Transaction) Identify() pb.Id {
	return pb.Id_TransactionType
}
