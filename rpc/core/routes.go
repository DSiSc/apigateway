package core

import (
	rpc "github.com/DSiSc/apigateway/rpc/lib/server"
)

// NOTE: Amino is registered in rpc/core/types/wire.go.
var Routes = map[string]*rpc.RPCFunc{

	// namespace "eth" API
	"eth_sendTransaction": rpc.NewRPCFunc(SendTransaction, "args"),
}
