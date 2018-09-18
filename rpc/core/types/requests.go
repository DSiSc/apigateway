package core_types

import (
	"github.com/DSiSc/apigateway/common"
	"github.com/DSiSc/apigateway/core/types"
)

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     types.Address  `json:"from"`
	To       *types.Address `json:"to"`
	Gas      *common.Uint64 `json:"gas"`
	GasPrice *common.Big    `json:"gasPrice"`
	Value    *common.Big    `json:"value"`
	Nonce    *common.Uint64 `json:"nonce,omitempty"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *common.Bytes `json:"data"`
	Input *common.Bytes `json:"input"`
}

type StringArgs struct {
	From string `json:"from"`
}
