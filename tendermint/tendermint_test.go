package tendermint

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestInitialize(t *testing.T) {
	app, err := Initialize("tcp://0.0.0.0:26658")
	if err != nil {
		t.Fatal(err)
	}
	app.Wait()
}

func TestReq(t *testing.T) {
	app, err := Initialize("tcp://0.0.0.0:26658")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := app.Client.DeliverTxSync([]byte("pct"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("%+v", resp))
	fmt.Println(app.Info(types.RequestInfo{Version: "1"}))
}

func jString(v interface{}) string {
	d, _ := json.Marshal(v)
	return string(d)
}
