package hash

import (
	"errors"
	"golang.org/x/crypto/sha3"
)

func SingleHash(b []byte) []byte {
	return hash(b)
}

func DoubleHash(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, errors.New("len of data is zero")
	}
	temp := hash(b)
	return hash(temp[:]), nil
}

func hash(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	_, _ = h.Write(data)
	return h.Sum(nil)
}
