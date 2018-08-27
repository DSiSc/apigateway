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
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
)

// ------------------------
// package Consts, Vars

const (
	badNibble = ^uint64(0)
	uintBits  = 32 << (uint64(^uint(0)) >> 63)
)

var (
	bigWordNibbles int
	bigT           = reflect.TypeOf((*Big)(nil))
	uintT          = reflect.TypeOf(Uint(0))
	uint64T        = reflect.TypeOf(Uint64(0))

	// errors
	ErrBig256Range = NewError("hex number > 256 bits")
	ErrUint64Range = NewError("hex number > 64 bits")
	ErrUintRange   = NewError(fmt.Sprintf("hex number > %d bits", uintBits))
)

// ------------------------
// package init

func init() {
	// This is a weird way to compute the number of nibbles required for big.Word.
	// The usual way would be to use constant arithmetic but go vet can't handle that.
	b, _ := new(big.Int).SetString("FFFFFFFFFF", 16)
	switch len(b.Bits()) {
	case 1:
		bigWordNibbles = 16
	case 2:
		bigWordNibbles = 8
	default:
		panic("weird big.Word size")
	}
}

// ------------------------
// package Struct Big

// Big marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
//
// Negative integers are not supported at this time. Attempting to marshal them will
// return an error. Values larger than 256bits are rejected by Unmarshal but will be
// marshaled without error.
type Big big.Int

// NewBig new Big from big.Int
func NewBig(bigint *big.Int) *Big {
    return (*Big)(bigint)
}

// MarshalText implements encoding.TextMarshaler
func (b Big) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Big) UnmarshalJSON(input []byte) error {
	if !IsString(input) {
		return errNonString(bigT)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), bigT)
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *Big) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 64 {
		return ErrBig256Range
	}
	words := make([]big.Word, len(raw)/bigWordNibbles+1)
	end := len(raw)
	for i := range words {
		start := end - bigWordNibbles
		if start < 0 {
			start = 0
		}
		for ri := start; ri < end; ri++ {
			nib := decodeNibble(raw[ri])
			if nib == badNibble {
				return ErrSyntax
			}
			words[i] *= 16
			words[i] += big.Word(nib)
		}
		end = start
	}
	var dec big.Int
	dec.SetBits(words)
	*b = (Big)(dec)
	return nil
}

// String returns the hex encoding of b.
func (b Big) String() string {

	nbits := b.ToBigInt().BitLen()
	if nbits == 0 {
		return "0x0"
	}
	return fmt.Sprintf("%#x", b.ToBigInt())
}

// ToBytes return []byte
func (b Big) ToBytes() []byte {
    return []byte(b.String())
}

// ToBigInt return big.Int
func (b Big) ToBigInt() *big.Int {
    return (*big.Int)(&b)
}

// -------------------------
// package Struct Unit64

// Uint64 marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type Uint64 uint64

// NewUint64 create Uint64 from uint64
func NewUint64(i uint64) *Uint64 {
    return (*Uint64)(&i)
}

// Unit64Encode encodes i as a hex string with 0x prefix.
func Unit64Encode(i uint64) string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, i, 16))
}

// MarshalText implements encoding.TextMarshaler.
func (b Uint64) MarshalText() ([]byte, error) {
    return []byte(Unit64Encode(uint64(b))), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint64) UnmarshalJSON(input []byte) error {
	if !IsString(input) {
		return errNonString(uint64T)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), uint64T)
}

// UnmarshalText implements encoding.TextUnmarshaler
func (b *Uint64) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 16 {
		return ErrUint64Range
	}
	var dec uint64
	for _, byte := range raw {
		nib := decodeNibble(byte)
		if nib == badNibble {
			return ErrSyntax
		}
		dec *= 16
		dec += nib
	}
	*b = Uint64(dec)
	return nil
}

// String returns the hex encoding of b.
func (b Uint64) String() string {
	return Unit64Encode(uint64(b))
}

// ToBytes return []byte
func (b Uint64) ToBytes() []byte {
    return []byte(b.String())
}

// Touint64 return uint64
func (b Uint64) Touint64() uint64 {
    return (uint64)(b)
}

// ------------------------
// package Struct Uint

// Uint marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type Uint uint

func NewUint(i uint) *Uint {
    return (*Uint)(&i)
}

// MarshalText implements encoding.TextMarshaler.
func (b Uint) MarshalText() ([]byte, error) {
	return Uint64(b).MarshalText()
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint) UnmarshalJSON(input []byte) error {
	if !IsString(input) {
		return errNonString(uintT)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), uintT)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Uint) UnmarshalText(input []byte) error {
	var u64 Uint64
	err := u64.UnmarshalText(input)
	if u64 > Uint64(^uint(0)) || err == ErrUint64Range {
		return ErrUintRange
	} else if err != nil {
		return err
	}
	*b = Uint(u64)
	return nil
}

// String returns the hex encoding of b.
func (b Uint) String() string {
	return Unit64Encode(uint64(b))
}

// ToBytes return the []byte
func (b Uint) ToBytes() []byte{
    return []byte(b.String())
}

// Touint return the uint
func (b Uint) Touint() uint {
    return (uint)(b)
}

// ---------------------------------
// package Function inner

func checkNumberText(input []byte) (raw []byte, err error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if !HexHasPrefix(string(input)) {
		return nil, ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return nil, ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return nil, ErrLeadingZero
	}
	return input, nil
}

func wrapTypeError(err error, typ reflect.Type) error {
	if _, ok := err.(*cmnError); ok {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: typ}
	}
	return err
}

func errNonString(typ reflect.Type) error {
	return &json.UnmarshalTypeError{Value: "non-string", Type: typ}
}

func decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return badNibble
	}
}
