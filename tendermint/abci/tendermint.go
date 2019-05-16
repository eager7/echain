package abci

import (
	"errors"
	"fmt"
	"github.com/eager7/elog"
	"github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	tlog "github.com/tendermint/tendermint/libs/log"
	"os"
)

var log = elog.NewLogger("tm", elog.InfoLevel)

type EChainApp struct {
	types.Application
	Server  common.Service
	Client  abcicli.Client
	AppHash []byte
}

func Initialize(sock string) (*EChainApp, error) {
	app := EChainApp{}
	ser := server.NewSocketServer(sock, &app) //fmt.Sprintf("unix://%s.sock", sock)
	if err := ser.Start(); err != nil {
		return nil, errors.New("start server err:" + err.Error())
	}
	fmt.Println("server status:", ser.String(), ser.IsRunning())

	cli := abcicli.NewSocketClient(sock, true)
	if err := cli.Start(); err != nil {
		_ = ser.Stop()
		return nil, errors.New("start client err:" + err.Error())
	}
	fmt.Println("client status:", cli.String(), cli.IsRunning())

	app.Server = ser
	app.Client = cli
	app.Wait()
	return &app, nil
}

func (e *EChainApp) Wait() {
	common.TrapSignal(tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout)), func() {
		if err := e.Server.Stop(); err != nil {
			fmt.Println("stop server err:", err)
		}
	})
}

// Info/Query Connection
func (e *EChainApp) Info(req types.RequestInfo) types.ResponseInfo {
	log.Info("app info:", req)
	return types.ResponseInfo{
		Data:             "echain abci",
		Version:          "1.0",
		AppVersion:       1,
		LastBlockHeight:  0,
		LastBlockAppHash: []byte("hash"),
	}
} // Return application info

func (e *EChainApp) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	log.Info("app set option:", req)
	return types.ResponseSetOption{}
} // Set application option

func (e *EChainApp) Query(req types.RequestQuery) types.ResponseQuery {
	log.Info("app query:", req.String())
	return types.ResponseQuery{
		Code:      types.CodeTypeOK,
		Log:       "query" + req.String(),
		Info:      "",
		Index:     0,
		Key:       nil,
		Value:     nil,
		Proof:     nil,
		Height:    0,
		Codespace: "",
	}
} // Query for state

func (e *EChainApp) CheckTx(tx []byte) types.ResponseCheckTx {
	log.Info("app check tx:", string(tx))
	return types.ResponseCheckTx{
		Code:                 0,
		Data:                 nil,
		Log:                  "",
		Info:                 "echain tx check",
		GasWanted:            0,
		GasUsed:              0,
		Tags:                 nil,
		Codespace:            "",
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
} // Validate a tx for the mempool

func (e *EChainApp) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	log.Info("app init chain:", req)
	return types.ResponseInitChain{}
} // Initialize blockchain with validators and other info from TendermintCore

func (e *EChainApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	log.Info("app begin block:", req.Header.NumTxs, req.Header.Height, req.Header.AppHash)
	e.AppHash = req.Header.AppHash
	return types.ResponseBeginBlock{}
} // Signals the beginning of a block

func (e *EChainApp) DeliverTx(tx []byte) types.ResponseDeliverTx {
	log.Info("app deliver tx:", string(tx))
	return types.ResponseDeliverTx{}
} // Deliver a tx for full processing

func (e *EChainApp) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	log.Info("app end block:", req.String())
	return types.ResponseEndBlock{}
} // Signals the end of a block, returns changes to the validator set

func (e *EChainApp) Commit() types.ResponseCommit {
	log.Info("app commit")
	return types.ResponseCommit{
		Data: e.AppHash,
	}
} // Commit the state and return the application Merkle root hash
