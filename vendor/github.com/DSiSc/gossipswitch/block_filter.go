package gossipswitch

import (
	"errors"
	"github.com/DSiSc/craft/types"
)

// TxFilter is an implemention of switch message filter,
// switch will use transaction filter to verify transaction message.
type BlockFilter struct {
}

// create a new block filter instance.
func NewBlockFilter() *BlockFilter {
	return &BlockFilter{}
}

// Verify verify a switch message whether is validated.
// return nil if message is validated, otherwise return relative error
func (blockValidator *BlockFilter) Verify(msg SwitchMsg) error {
	switch msg := msg.(type) {
	case *types.Block:
		return blockValidator.doValidate(msg)
	default:
		return errors.New("unsupported message type")
	}
}

// do verify operation
func (blockValidator *BlockFilter) doValidate(block *types.Block) error {
	//TODO
	return nil
}
