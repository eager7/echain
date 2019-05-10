package types

import (
	"errors"
	"github.com/eager7/echain/core/bloom"
	"github.com/eager7/echain/core/pb"
	"github.com/golang/protobuf/proto"
)

type Header struct {
	Version     uint32      `json:"version"`
	ChainID     Hash        `json:"chain_id"`
	TimeStamp   int64       `json:"time_stamp"`
	Height      uint64      `json:"height"`
	PrevHash    Hash        `json:"prev_hash"`
	MerkleHash  Hash        `json:"merkle_hash"`
	StateHash   Hash        `json:"state_hash"`
	ReceiptHash Hash        `json:"receipt_hash"`
	Bloom       bloom.Bloom `json:"bloom"`
	Signatures  Signature   `json:"signatures"`
	Hash        Hash        `json:"hash"`
}

func (h *Header) Serialize() ([]byte, error) {
	d, err := proto.Marshal(&pb.Header{
		Version:     h.Version,
		ChainID:     h.ChainID.Bytes(),
		Timestamp:   h.TimeStamp,
		Height:      h.Height,
		PrevHash:    h.PrevHash.Bytes(),
		MerkleHash:  h.MerkleHash.Bytes(),
		StateHash:   h.StateHash.Bytes(),
		ReceiptHash: h.ReceiptHash.Bytes(),
		Bloom:       h.Bloom.Bytes(),
		Hash:        h.Hash.Bytes(),
		Sign: &pb.Signature{
			PubKey:  h.Signatures.PubKey,
			SigData: h.Signatures.SigData,
		},
	})
	if err != nil {
		return nil, errors.New("marshal err:" + err.Error())
	}
	return d, nil
}

func (h *Header) Deserialize(data []byte) error {

	return nil
}

func (h *Header) Instance() interface{} {
	return h
}

func (h *Header) Identify() pb.Id {
	return pb.Id_HeaderType
}
