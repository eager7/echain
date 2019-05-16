package mpt_test

import (
	"fmt"
	"github.com/eager7/echain/common"
	"github.com/eager7/echain/core/etypes"
	mpt "github.com/eager7/echain/core/state"
	"os"
	"testing"
)

func TestNewMptTree(t *testing.T) {
	_ = os.RemoveAll("/tmp/tree")
	fmt.Println(mpt.NewMptTree("/tmp/tree", etypes.HashSetBytes(common.SingleHash([]byte("test")))))
}
