package tclient

import (
	"fmt"
	"github.com/eager7/echain/common"
	"github.com/eager7/echain/core/etypes"
	"math/big"
	"testing"
	"time"
)

func TestSendTx(t *testing.T) {
	cli := Initialize("tcp://0.0.0.0:26657")
	h, _ := common.DoubleHash([]byte("pct"))
	tx := etypes.Transaction{
		Version:    1,
		ChainID:    etypes.HashSetBytes(h),
		From:       etypes.AddressSetBytes([]byte("form")),
		To:         etypes.AddressSetBytes([]byte("to")),
		Amount:     new(big.Int).SetUint64(100),
		TimeStamp:  time.Now().UnixNano(),
		Payload:    nil,
		Signatures: etypes.Signature{},
		Hash:       etypes.HashSetBytes(h),
	}
	_, err := cli.SendTx(&tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryBlock(t *testing.T) {
	cli := Initialize("tcp://0.0.0.0:26657")
	block, txs, err := cli.QueryBlock(1915)
	if err != nil {
		t.Fatal(err)
	}
	for _, tx := range txs {
		fmt.Printf("tx:%+v\n", tx)
	}
	fmt.Println(block.String())
}
