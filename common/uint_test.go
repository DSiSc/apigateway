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
	"testing"

	inl "github.com/DSiSc/apigateway/common/internal"
)

// ---------------------------
// package Struct

// ----------------------------
// package Consts, Vars
var (
	encodeUintTests = []marshalTest{
		{uint(0), "0x0"},
		{uint(1), "0x1"},
		{uint(0xff), "0xff"},
		{uint(0x11223344), "0x11223344"},
	}

	unmarshalUintTests = []unmarshalTest{
		// invalid encoding
		{input: "", wantErr: errJSONEOF},
		{input: "null", wantErr: errNonString(uintT)},
		{input: "10", wantErr: errNonString(uintT)},
		{input: `"0"`, wantErr: wrapTypeError(ErrMissingPrefix, uintT)},
		{input: `"0x"`, wantErr: wrapTypeError(ErrEmptyNumber, uintT)},
		{input: `"0x01"`, wantErr: wrapTypeError(ErrLeadingZero, uintT)},
		{input: `"0x100000000"`, want: uint(maxUint33bits), wantErr32bit: wrapTypeError(ErrUintRange, uintT)},
		{input: `"0xfffffffffffffffff"`, wantErr: wrapTypeError(ErrUintRange, uintT)},
		{input: `"0xx"`, wantErr: wrapTypeError(ErrSyntax, uintT)},
		{input: `"0x1zz01"`, wantErr: wrapTypeError(ErrSyntax, uintT)},

		// valid encoding
		{input: `""`, want: uint(0)},
		{input: `"0x0"`, want: uint(0)},
		{input: `"0x2"`, want: uint(0x2)},
		{input: `"0x2F2"`, want: uint(0x2f2)},
		{input: `"0X2F2"`, want: uint(0x2f2)},
		{input: `"0x1122aaff"`, want: uint(0x1122aaff)},
		{input: `"0xbbb"`, want: uint(0xbbb)},
		{input: `"0xffffffff"`, want: uint(0xffffffff)},
		{input: `"0xffffffffffffffff"`, want: uint(maxUint64bits), wantErr32bit: wrapTypeError(ErrUintRange, uintT)},
	}
)

// --------------------------
// package Test* json encoding
func TestMarshalUint(t *testing.T) {
	for _, test := range encodeUintTests {
		in := test.input.(uint)
		out, err := json.Marshal(Uint(in))
		if err != nil {
			t.Errorf("%d: %v", in, err)
			continue
		}
		if want := `"` + test.want + `"`; string(out) != want {
			t.Errorf("%d: MarshalJSON output mismatch: got %q, want %q", in, out, want)
			continue
		}
		if out := (*Uint)(&in).String(); out != test.want {
			t.Errorf("%x: String mismatch: got %q, want %q", in, out, test.want)
			continue
		}
	}
}

func TestUnmarshalUint(t *testing.T) {
	for _, test := range unmarshalUintTests {
		var v Uint
		err := json.Unmarshal([]byte(test.input), &v)
		if uintBits == 32 && test.wantErr32bit != nil {
			inl.CheckError(t, test.input, err, test.wantErr32bit)
			continue
		}
		if !inl.CheckError(t, test.input, err, test.wantErr) {
			continue
		}
		if uint(v) != test.want.(uint) {
			t.Errorf("input %s: value mismatch: got %d, want %d", test.input, v, test.want)
			continue
		}
	}
}

func TestToUint(t *testing.T) {

	input := uint(1)
	want := NewUint(input).Touint()

	if input != want {
		t.Errorf("Error with to unit: got %d, want %d", input, want)
	}
}

// -------------------------
// package Functions inner
