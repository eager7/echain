package etypes

import (
	"errors"
	"github.com/eager7/echain/core/pb"
	"github.com/golang/protobuf/proto"
)

type Signature struct {
	PubKey  []byte `json:"pub_key"`
	SigData []byte `json:"sig_data"`
}

func (s *Signature) Serialize() ([]byte, error) {
	d, err := proto.Marshal(&pb.Signature{SigData: s.SigData, PubKey: s.PubKey})
	if err != nil {
		return nil, errors.New("marshal err:" + err.Error())
	}
	return d, nil
}

func (s *Signature) Deserialize(data []byte) error {
	if len(data) == 0 {
		return errors.New("input data's length is zero")
	}
	var sig pb.Signature
	if err := proto.Unmarshal(data, &sig); err != nil {
		return errors.New("unmarshal err:" + err.Error())
	}
	s.PubKey = sig.PubKey
	s.SigData = sig.SigData
	return nil
}

func (s *Signature) Instance() interface{} {
	return s
}

func (s *Signature) Identify() pb.Id {
	return pb.Id_SignatureType
}
