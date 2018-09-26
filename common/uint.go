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
	inl "github.com/DSiSc/apigateway/common/internal"
)

// ------------------------
// package Consts, Vars

// ------------------------
// package Struct Uint

// Uint marshals/unmarshals as a JSON string with 0x prefix.
type Uint uint

func NewUint(i uint) *Uint {
	return (*Uint)(&i)
}

// MarshalText implements encoding.TextMarshaler.
func (b Uint) MarshalText() ([]byte, error) {
	return Uint64(b).MarshalText()
}

// MarshalJSON implements encoding.JSONMarshaler.
func (b Uint) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%v", b))
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint) UnmarshalJSON(input []byte) error {
	if !inl.IsString(input) {
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
	return string(Ghex.EncodeUint(b.Touint()))
}

// ToBytes return the []byte
func (b Uint) ToBytes() []byte {
	return []byte(b.String())
}

// Touint return the uint
func (b Uint) Touint() uint {
	return (uint)(b)
}
