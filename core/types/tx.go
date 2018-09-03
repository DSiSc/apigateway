// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"math/big"

	"crypto/sha256"
	"encoding/json"

	"github.com/DSiSc/craft/types"
)

func Sum(bz []byte) []byte {
	hash := sha256.Sum256(bz)
	return hash[:types.HashLength]
}

func TxHash(tx *types.Transaction) (hash types.Hash) {
	jsonByte, _ := json.Marshal(tx)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

func HashBytes(a types.Hash) []byte {
	b := make([]byte, len(a))
	copy(b, a[:])
	return b
}

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

func TypeConvert(a *Address) *types.Address {
	var address types.Address
	copy(address[:], a[:])
	return &address
}

// New a transaction
func newTransaction(nonce uint64, to *Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, from *Address) *types.Transaction {
	if len(data) > 0 {
		data = CopyBytes(data)
	}
	d := types.TxData{
		AccountNonce: nonce,
		Recipient:    TypeConvert(to),
		From:         TypeConvert(from),
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     gasLimit,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &types.Transaction{Data: d}
}

func NewTransaction(nonce uint64, to Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, from Address) *types.Transaction {
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data, &from)
}
