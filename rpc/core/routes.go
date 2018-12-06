package core

import (
	rpc "github.com/DSiSc/apigateway/rpc/lib/server"
)

// NOTE: Amino is registered in rpc/core/types/wire.go.
var Routes = map[string]*rpc.RPCFunc{
	// namespace "eth" API
	"eth_sendTransaction":                     rpc.NewRPCFunc(SendTransaction, "args"),
	"eth_getBlockByHash":                      rpc.NewRPCFunc(GetBlockByHash, "blockHash, fullTx"),
	"eth_getBlockByNumber":                    rpc.NewRPCFunc(GetBlockByNumber, "blockNr, fullTx"),
	"eth_getTransactionByHash":                rpc.NewRPCFunc(GetTransactionByHash, "hash"),
	"eth_getTransactionReceipt":               rpc.NewRPCFunc(GetTransactionReceipt, "hash"),
	"eth_getBlockTransactionCountByHash":      rpc.NewRPCFunc(GetBlockTransactionCountByHash, "blockHash"),
	"eth_getBlockTransactionCountByNumber":    rpc.NewRPCFunc(GetBlockTransactionCountByNumber, "blockNr"),
	"eth_blockNumber":                         rpc.NewRPCFunc(BlockNumber, ""),
	"eth_getBalance":                          rpc.NewRPCFunc(GetBalance, "address, blockNr"),
	"eth_getCode":                             rpc.NewRPCFunc(GetCode, "address, blockNr"),
	"eth_getTransactionCount":                 rpc.NewRPCFunc(GetTransactionCount, "address, blockNr"),
	"eth_getTransactionByBlockHashAndIndex":   rpc.NewRPCFunc(GetTransactionByBlockHashAndIndex, "blockHash, index"),
	"eth_getTransactionByBlockNumberAndIndex": rpc.NewRPCFunc(GetTransactionByBlockNumberAndIndex, "blockNr, index"),
	"eth_call":        rpc.NewRPCFunc(Call, "args, blockNr"),
	"eth_gasPrice":    rpc.NewRPCFunc(GasPrice, ""),
	"eth_estimateGas": rpc.NewRPCFunc(EstimateGas, "args"),
	"eth_accounts":    rpc.NewRPCFunc(Accounts, ""),
	"net_listening":   rpc.NewRPCFunc(Listening, ""),
	"net_version":     rpc.NewRPCFunc(Version, ""),
	"net_nodeInfo":    rpc.NewRPCFunc(NodeInfo, ""),
	"net_channelInfo": rpc.NewRPCFunc(ChannelInfo, ""),
}

func AddTestRoutes() {
	Routes["echo"] = rpc.NewRPCFunc(EchoResult, "arg")
	Routes["echo_args"] = rpc.NewRPCFunc(EchoResultArgs, "arg")
}
