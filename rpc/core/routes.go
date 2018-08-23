package core

import (
	rpc "github.com/DSiSc/apigateway/rpc/lib/server"
)

// NOTE: Amino is registered in rpc/core/types/wire.go.
var Routes = map[string]*rpc.RPCFunc{
	// namespace "eth" API
	"eth_sendTransaction": rpc.NewRPCFunc(SendTransaction, "args"),
}

func AddTestRoutes() {
	Routes["echo"] = rpc.NewRPCFunc(EchoResult, "arg")
	Routes["echo_args"] = rpc.NewRPCFunc(EchoResultArgs, "arg")
}
