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
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	cmn "github.com/DSiSc/apigateway/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

// Lengths of hashes and addresses in bytes.
const (
	AddressLength = 20
)

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address {
	return BytesToAddress(b.Bytes())
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address {
	return BytesToAddress(cmn.Ghex.MustDecodeString(s))
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalJSON() ([]byte, error) {
	js := strconv.Quote(cmn.Ghex.EncodeToString(a[:]))
	return []byte(js), nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {

	unQuote, err := strconv.Unquote(string(input))
	if err != nil {
		return err

	}

	srcBytes := []byte(unQuote)
	dstBytes := make([]byte, cmn.Ghex.DecodeLen(len(srcBytes)))
	_, errDecode := cmn.Ghex.Decode(dstBytes, srcBytes)
	if errDecode != nil {
		return errDecode
	}

	a.SetBytes(dstBytes)
	return nil
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Big converts an address to a big integer.
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() cmn.Hash { return cmn.BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}
