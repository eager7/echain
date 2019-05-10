package bloom

import (
	"math/big"
	"testing"
)

func TestNewBloom(t *testing.T) {
	b := NewBloom(nil)
	b.Add([]byte("pct"))
	if b.Test([]byte("pct")) != true {
		t.Fatal("error test pct")
	}
	if b.Test([]byte("pct2")) == true {
		t.Fatal("error test pct2")
	}

	data := b.Bytes()
	b2 := NewBloom(data)
	if b2.Test([]byte("pct")) != true {
		t.Fatal("error test pct")
	}
	if b2.Test([]byte("pct2")) == true {
		t.Fatal("error test pct2")
	}
}

func TestBloomCycle(t *testing.T) {
	b := NewBloom(nil)
	for i := 0; i < 100000; i++ {
		key := new(big.Int).SetInt64(int64(i))
		b.Add(key.Bytes())
	}
	for i := 0; i < 100000; i++ {
		key := new(big.Int).SetInt64(int64(i))
		if true != b.Test(key.Bytes()) {
			t.Fatal("test error", i)
		}
	}
}
