package core

import (
	rpc "github.com/DSiSc/apigateway/rpc/lib/server"
)

// NOTE: Amino is registered in rpc/core/types/wire.go.
var Routes = map[string]*rpc.RPCFunc{
	// namespace "eth" API
	"eth_sendTransaction": rpc.NewRPCFunc(SendTransaction, "args"),
	"eth_getBlockByHash": rpc.NewRPCFunc(GetBlockByHash, "blockHash, fullTx"),
	"eth_getBlockByNumber": rpc.NewRPCFunc(GetBlockByNumber, "blockNr, fullTx"),
	"eth_getTransactionByHash": rpc.NewRPCFunc(GetTransactionByHash, "hash"),
	"eth_getTransactionReceipt": rpc.NewRPCFunc(GetTransactionReceipt, "hash"),
}

func AddTestRoutes() {
	Routes["echo"] = rpc.NewRPCFunc(EchoResult, "arg")
	Routes["echo_args"] = rpc.NewRPCFunc(EchoResultArgs, "arg")
}
