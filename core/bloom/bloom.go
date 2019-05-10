package bloom

import (
	"fmt"
	"github.com/eager7/echain/common/hash"
	"math/big"
)

const (
	BloomByteLength = 1024
	BloomHashCycle  = 2
	BloomBitLength  = 8 * BloomByteLength
)

type Bloom [BloomByteLength]byte

func NewBloom(b []byte) Bloom {
	var bloom Bloom
	bloom.SetBytes(b)
	return bloom
}
func (b *Bloom) SetBytes(d []byte) {
	if len(b) < len(d) {
		panic(fmt.Sprintf("bloom bytes too big %d %d", len(b), len(d)))
	}
	copy(b[BloomByteLength-len(d):], d)
}

func (b Bloom) Bytes() []byte {
	return b[:]
}

func (b *Bloom) Add(data []byte) {
	b.add(new(big.Int).SetBytes([]byte(data)))
}

func (b *Bloom) add(d *big.Int) {
	bin := new(big.Int).SetBytes(b[:])
	bin.Or(bin, bloom(d.Bytes()))
	b.SetBytes(bin.Bytes())
}

func (b Bloom) Test(test []byte) bool {
	return b.test(new(big.Int).SetBytes(test))
}

func (b Bloom) test(test *big.Int) bool {
	return bloomLookup(b, test.Bytes())
}

func (b Bloom) Big() *big.Int {
	return new(big.Int).SetBytes(b[:])
}

func bloom(b []byte) *big.Int {
	b = hash.SingleHash(b[:])
	r := new(big.Int)
	for i := 0; i < BloomHashCycle; i += 2 {
		t := big.NewInt(1)
		b := (uint(b[i+1]) + (uint(b[i]) << 8)) & 2047
		r.Or(r, t.Lsh(t, b))
	}
	return r
}
func bloomLookup(bin Bloom, key []byte) bool {
	b := bin.Big()
	cmp := bloom(key)
	return b.And(b, cmp).Cmp(cmp) == 0
}
