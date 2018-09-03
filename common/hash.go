// Copyright(c) 2018 DSiSc Group All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
)

// -------------------------
// package Consts, Vars
// Lengths of hashes and addresses in bytes.
const (
	HashLength = 32
)

// ------------------
// package Functions

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash {
	return BytesToHash(b.Bytes())
}

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash {
	return BytesToHash(Ghex.MustDecodeString(s))
}

// ----------------------
// package Struct Hash

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

func (h Hash) MarshalJSON() ([]byte, error) {
	js := strconv.Quote(Ghex.EncodeToString(h[:]))
	return []byte(js), nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {

	unQuote, err := strconv.Unquote(string(input))
	if err != nil {
		return err

	}

	srcBytes := []byte(unQuote)
	dstBytes := make([]byte, Ghex.DecodeLen(len(srcBytes)))
	_, errDecode := Ghex.Decode(dstBytes, srcBytes)
	if errDecode != nil {
		return errDecode
	}

	h.SetBytes(dstBytes)
	return nil
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string {
	return Ghex.EncodeToString(h[:])
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (h Hash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h[:])
}
