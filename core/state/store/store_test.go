package store_test

import (
	"bytes"
	"fmt"
	"github.com/eager7/echain/core/state/store"
	"os"
	"testing"
)

var s store.Storage

func TestStore(t *testing.T) {
	s, err := store.NewBlockStore("/tmp/test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s.SearchAll())
	fmt.Println(s.Get([]byte("test")))
	//存储
	_ = s.Put([]byte("key1"), []byte("value1"))
	_ = s.Put([]byte("key2"), []byte("value2"))
	_ = s.Put([]byte("key3"), []byte("value3"))
	//读取
	if v, err := s.Get([]byte("key1")); err != nil {
		if !bytes.Equal(v, []byte("value1")) {
			t.Fatal("value1 error")
		}
	}
	if v, err := s.Get([]byte("key3")); err != nil {
		if !bytes.Equal(v, []byte("value3")) {
			t.Fatal("value1 error")
		}
	}
	//批处理存储
	s.BatchPut([]byte("key01"), []byte("value01"))
	s.BatchPut([]byte("key02"), []byte("value02"))
	s.BatchPut([]byte("key03"), []byte("value03"))
	s.BatchPut([]byte("key04"), []byte("value04"))
	_ = s.BatchCommit()
	if v, err := s.Get([]byte("key03")); err != nil {
		if !bytes.Equal(v, []byte("value03")) {
			t.Fatal("value1 error")
		}
	}
	//迭代器遍历
	it := s.NewIterator()
	for it.Next() {
		fmt.Println(string(it.Key()), string(it.Value()))
	}
	it.Release()
	if err := it.Error(); err != nil {
		t.Fatal(err)
	}

	//函数遍历
	re, err := s.SearchAll()
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range re {
		fmt.Println(k, v)
	}
}

func BenchmarkLevelDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = s.Put([]byte("key"), []byte("value"))
		v, _ := s.Get([]byte("key"))
		fmt.Println(v)
	}
}

func init() {
	_ = os.RemoveAll("/tmp/store_benchmark")
	s, _ = store.NewBlockStore("/tmp/store_benchmark")
}
