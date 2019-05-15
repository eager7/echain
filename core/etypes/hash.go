package etypes

import (
	"bytes"
	"encoding/hex"
)

const HashLen = 32

type Hash [HashLen]byte

func HashSetBytes(b []byte) Hash {
	var hash Hash
	if len(b) > len(hash) {
		b = b[len(b)-HashLen:]
	}
	copy(hash[HashLen-len(b):], b)
	return hash
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) Hex() string {
	s := hex.EncodeToString(h[:])
	if len(s) == 0 {
		s = "0"
	}
	return "0x" + s
}

func HashSetHex(data string) (hash Hash) {
	if len(data) > 1 {
		if data[0:2] == "0x" || data[0:2] == "0X" {
			data = data[2:]
		}
	}
	if len(data)%2 == 1 {
		data = "0" + data
	}
	h, _ := hex.DecodeString(data)
	return HashSetBytes(h)
}

func (h *Hash) Equals(b *Hash) bool {
	if nil == h {
		return nil == b
	}
	if nil == b {
		return false
	}
	return bytes.Equal(h[:], b[:])
}