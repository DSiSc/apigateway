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

	gconf "github.com/DSiSc/craft/config"
	"github.com/DSiSc/craft/rlp"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto/sha3"
	"hash"
)

const (
	DefaultGasPrice = 1
)

// get hash algorithm by global config
func HashAlg() hash.Hash {
	var alg string
	if value, ok := gconf.GlobalConfig.Load(gconf.HashAlgName); ok {
		alg = value.(string)
	} else {
		alg = "SHA256"
	}
	return sha3.NewHashByAlgName(alg)
}

// calculate the hash value of the rlp encoded byte of x
func rlpHash(x interface{}) (h types.Hash) {
	hw := HashAlg()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// TxHash calculate tx's hash
func TxHash(tx *types.Transaction) (hash types.Hash) {
	if hash := tx.Hash.Load(); hash != nil {
		return hash.(types.Hash)
	}
	v := rlpHash(tx)
	tx.Hash.Store(v)
	return v
}

// HeaderHash calculate block's hash
func HeaderHash(block *types.Block) (hash types.Hash) {
	//var defaultHash types.Hash
	if !(block.HeaderHash == types.Hash{}) {
		var hash types.Hash
		copy(hash[:], block.HeaderHash[:])
		return hash
	}
	return rlpHash(block.Header)
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
	if a != nil {
		copy(address[:], a[:])
		return &address
	}
	return nil
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

func NewTransaction(nonce uint64, to *Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, from Address) *types.Transaction {
	return newTransaction(nonce, to, amount, gasLimit, gasPrice, data, &from)
}
