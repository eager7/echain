package types

import (
	"errors"
	"github.com/eager7/echain/core/pb"
	"math/big"
)

type Receipt struct {
	From      Address  `json:"from"`
	To        Address  `json:"to"`
	Amount    *big.Int `json:"amount"`
	Timestamp int64    `json:"timestamp"`
	Result    []byte   `json:"result"`
	Hash      Hash     `json:"hash"`
}

func (r *Receipt) Serialize() ([]byte, error) {
	d, err := proto.Marshal(&pb.Receipt{})
	if err != nil {
		return nil, errors.New("marshal err:" + err.Error())
	}
	return d, nil
}

func (r *Receipt) Deserialize(data []byte) error {

	return nil
}

func (r *Receipt) Instance() interface{} {
	return r
}

func (r *Receipt) Identify() pb.Id {
	return pb.Id_TransactionType
}
