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
	"strconv"
	"testing"

	inl "github.com/DSiSc/apigateway/common/internal"
)

// ---------------------------
// package Struct

// ----------------------------
// package Consts, Vars

// --------------------------
// package Test* json encoding
func TestUnmarshalUint64(t *testing.T) {
	var unmarshalUint64Tests = []unmarshalTest{
		// invalid encoding
		{input: "", wantErr: errJSONEOF},
		{input: "null", wantErr: strconv.ErrSyntax},
		{input: "10", wantErr: strconv.ErrSyntax},
		{input: `"0"`, wantErr: wrapTypeError(ErrMissingPrefix, uint64T)},
		{input: `"0x"`, wantErr: wrapTypeError(ErrEmptyNumber, uint64T)},
		{input: `"0x01"`, wantErr: wrapTypeError(ErrLeadingZero, uint64T)},
		{input: `"0xfffffffffffffffff"`, wantErr: wrapTypeError(ErrUint64Range, uint64T)},
		{input: `"0xx"`, wantErr: wrapTypeError(ErrSyntax, uint64T)},
		{input: `"0x1zz01"`, wantErr: wrapTypeError(ErrSyntax, uint64T)},

		// valid encoding
		{input: `""`, want: uint64(0)},
		{input: `"0x0"`, want: uint64(0)},
		{input: `"0x2"`, want: uint64(0x2)},
		{input: `"0x2F2"`, want: uint64(0x2f2)},
		{input: `"0X2F2"`, want: uint64(0x2f2)},
		{input: `"0x1122aaff"`, want: uint64(0x1122aaff)},
		{input: `"0xbbb"`, want: uint64(0xbbb)},
		{input: `"0xffffffffffffffff"`, want: uint64(0xffffffffffffffff)},
	}
	for _, test := range unmarshalUint64Tests {
		var v Uint64
		err := json.Unmarshal([]byte(test.input), &v)
		if !inl.CheckError(t, test.input, err, test.wantErr) {
			continue
		}
		if uint64(v) != test.want.(uint64) {
			t.Errorf("input %s: value mismatch: got %d, want %d", test.input, v, test.want)
			continue
		}
	}
}

func BenchmarkUnmarshalUint64(b *testing.B) {
	input := []byte(`"0x123456789abcdf"`)
	for i := 0; i < b.N; i++ {
		var v Uint64
		v.UnmarshalJSON(input)
	}
}

func TestMarshalUint64(t *testing.T) {

	var encodeUint64Tests = []marshalTest{
		{uint64(0), "0x0"},
		{uint64(1), "0x1"},
		{uint64(0xff), "0xff"},
		{uint64(0x1122334455667788), "0x1122334455667788"},
	}

	for _, test := range encodeUint64Tests {
		in := test.input.(uint64)
		out, err := json.Marshal(Uint64(in))
		if err != nil {
			t.Errorf("%d: %v", in, err)
			continue
		}
		if want := `"` + test.want + `"`; string(out) != want {
			t.Errorf("%d: MarshalJSON output mismatch: got %q, want %q", in, out, want)
			continue
		}
		if out := (*Uint64)(&in).String(); out != test.want {
			t.Errorf("%x: String mismatch: got %q, want %q", in, out, test.want)
			continue
		}
	}
}

// ----------------------
// package Test* others
func TestToUint64(t *testing.T) {

	input := uint64(1)
	want := NewUint64(input).Touint64()

	if input != want {
		t.Errorf("Error with to unit64: got %d, want %d", input, want)
	}
}

// -------------------------
// package Functions inner
