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
	"strconv"
)

// ---------------------------
// package Struct Bytes

type Bytes []byte

// New new Bytes
func NewBytes(data []byte) *Bytes {
	var bz Bytes
	bz = append(bz, data...)
	return &bz
}

// Marshal needed for protobuf compatibility
func (bz Bytes) Marshal() ([]byte, error) {
	return bz, nil
}

// Unmarshal needed for protobuf compatibility
func (bz *Bytes) Unmarshal(data []byte) error {
	*bz = data
	return nil
}

// MarshalJSON implement encoding/json Marshaler interface.
func (bz Bytes) MarshalJSON() ([]byte, error) {
	js := strconv.Quote(Ghex.EncodeToString(bz))
	return []byte(js), nil
}

// This is the point of Bytes.
func (bz *Bytes) UnmarshalJSON(data []byte) error {
	if string(data) == `""` {
		*bz = []byte(``)
		return nil
	}

	unQuote, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	bzDecode, err := Ghex.DecodeString(unQuote)
	if err != nil {
		return err
	}

	*bz = bzDecode
	return nil
}

func (bz *Bytes) UnmarshalText(data []byte) error {
	return bz.UnmarshalJSON(data)
}

// Allow it to fulfill various interfaces in light-client, etc...
func (bz Bytes) Bytes() []byte {
	return bz
}

func (bz Bytes) String() string {
	return Ghex.EncodeToString(bz)
}

func (bz Bytes) Format(s fmt.State, verb rune) {
	switch verb {
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", bz)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(bz))))
	}
}
