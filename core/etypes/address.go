package etypes

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/ripemd160"
)

const AddressLen = 20

type Address [AddressLen]byte

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Hex() string {
	h := hex.EncodeToString(a[:])
	if len(h) == 0 {
		h = "0"
	}
	return "0x" + h
}

func (a *Address) Equals(b *Address) bool {
	if nil == a {
		return nil == b
	}
	if nil == b {
		return false
	}
	return bytes.Equal(a[:], b[:])
}

func AddressSetHex(addr string) Address {
	if len(addr) > 1 {
		if addr[0:2] == "0x" || addr[0:2] == "0X" {
			addr = addr[2:]
		}
	}
	if len(addr)%2 == 1 {
		addr = "0" + addr
	}
	h, _ := hex.DecodeString(addr)
	return AddressSetBytes(h)
}

func AddressSetBytes(addr []byte) (address Address) {
	copy(address[:], addr)
	return address
}

func AddressSetPubKey(pubKey []byte) (address Address) {
	temp := sha256.Sum256(pubKey)
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(address[:0])
	address[0] = 0x01
	return address
}
