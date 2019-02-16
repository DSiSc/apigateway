package types

import (
	"encoding/json"
	"github.com/DSiSc/craft/types"
)

// HeaderHash calculate block's hash
func HeaderHash(block *types.Block) (hash types.Hash) {
	jsonByte, _ := json.Marshal(block.Header)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}
