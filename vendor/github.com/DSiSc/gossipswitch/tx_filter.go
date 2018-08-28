package gossipswitch

import (
	"errors"
	"github.com/DSiSc/craft/types"
)

// TxFilter is an implemention of switch message filter,
// switch will use transaction filter to verify transaction message.
type TxFilter struct {
}

// create a new transaction filter instance.
func NewTxFilter() *TxFilter {
	return &TxFilter{}
}

// Verify verify a switch message whether is validated.
// return nil if message is validated, otherwise return relative error
func (txValidator *TxFilter) Verify(msg SwitchMsg) error {
	switch msg := msg.(type) {
	case *types.Transaction:
		return txValidator.doVerify(msg)
	default:
		return errors.New("unsupported message type")
	}
}

// do verify operation
func (txValidator *TxFilter) doVerify(tx *types.Transaction) error {
	//TODO
	return nil
}
