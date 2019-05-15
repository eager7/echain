package main

import (
	"fmt"
	"github.com/eager7/echain/tendermint"
)

func main() {
	fmt.Println("start chain process...")
	app, err := tendermint.Initialize("tcp://0.0.0.0:26658")
	if err != nil {
		panic(err)
	}
	fmt.Println(app)
	select{}
}
