package etypes

import "github.com/eager7/echain/core/pb"

type EMessage interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	Instance() interface{}
	Identify() pb.Id
}
