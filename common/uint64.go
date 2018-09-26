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
	//inl "github.com/DSiSc/apigateway/common/internal"

	"encoding/json"
	"fmt"
	"strconv"
)

// ------------------------
// package Consts, Vars

// -------------------------
// package Struct Unit64

// Uint64 marshals/unmarshals as a JSON string with 0x prefix.
// The zero value marshals as "0x0".
type Uint64 uint64

// NewUint64 create Uint64 from uint64
func NewUint64(i uint64) *Uint64 {
	return (*Uint64)(&i)
}

// MarshalText implements encoding.TextMarshaler.
func (b Uint64) MarshalText() ([]byte, error) {
	return Ghex.EncodeUint64(b.Touint64()), nil
}

// MarshalJSON implements encoding.JSONMarshaler.
func (b Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%v", b))
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint64) UnmarshalJSON(input []byte) error {
	s, err := strconv.Unquote(string(input))
	if err != nil {
		return err
	}

	//	if !inl.IsString(input) {
	//return errNonString(uint64T)
	//}
	return wrapTypeError(b.UnmarshalText([]byte(s)), uint64T)
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
	return string(Ghex.EncodeUint64(b.Touint64()))
}

// ToBytes return []byte
func (b Uint64) ToBytes() []byte {
	return []byte(b.String())
}

// Touint64 return uint64
func (b Uint64) Touint64() uint64 {
	return (uint64)(b)
}
