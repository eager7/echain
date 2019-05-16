# echain
用Tendermint开发的一条测试公链

## 运行
首先启动Tendermint
```bash
tendermint node
```
然后启动main.go程序，此时会出空块，然后通过测试程序（tendermint/tclient）发送交易到Tendermint：
```bash
go test -v -run=TestSendTx
```
此时可以从Tendermint的打印中看到出块信息，用测试程序可以取到区块数据：
```bash
go test -v -run=TestQueryBlock
```