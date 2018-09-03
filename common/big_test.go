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
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	inl "github.com/DSiSc/apigateway/common/internal"
)

// ---------------------------
// package Struct

type unmarshalTest struct {
	input        string
	want         interface{}
	wantErr      error // if set, decoding must fail on any platform
	wantErr32bit error // if set, decoding must fail on 32bit platforms (used for Uint tests)
}

type marshalTest struct {
	input interface{}
	want  string
}

// ----------------------------
// package Consts, Vars
var (

	// These are variables (not constants) to avoid constant overflow
	// checks in the compiler on 32bit platforms.
	maxUint33bits = uint64(^uint32(0)) + 1
	maxUint64bits = ^uint64(0)

	// ---------------------
	// Test Cases
	encodeBigTests = []marshalTest{
		{referenceBig("0"), "0x0"},
		{referenceBig("1"), "0x1"},
		{referenceBig("ff"), "0xff"},
		{referenceBig("112233445566778899aabbccddeeff"), "0x112233445566778899aabbccddeeff"},
		{referenceBig("80a7f2c1bcc396c00"), "0x80a7f2c1bcc396c00"},
		{referenceBig("-80a7f2c1bcc396c00"), "-0x80a7f2c1bcc396c00"},
	}

	unmarshalBigTests = []unmarshalTest{
		// invalid encoding
		{input: "", wantErr: errJSONEOF},
		{input: "null", wantErr: errNonString(bigT)},
		{input: "10", wantErr: errNonString(bigT)},
		{input: `"0"`, wantErr: wrapTypeError(ErrMissingPrefix, bigT)},
		{input: `"0x"`, wantErr: wrapTypeError(ErrEmptyNumber, bigT)},
		{input: `"0x01"`, wantErr: wrapTypeError(ErrLeadingZero, bigT)},
		{input: `"0xx"`, wantErr: wrapTypeError(ErrSyntax, bigT)},
		{input: `"0x1zz01"`, wantErr: wrapTypeError(ErrSyntax, bigT)},
		{
			input:   `"0x10000000000000000000000000000000000000000000000000000000000000000"`,
			wantErr: wrapTypeError(ErrBig256Range, bigT),
		},

		// valid encoding
		{input: `""`, want: big.NewInt(0)},
		{input: `"0x0"`, want: big.NewInt(0)},
		{input: `"0x2"`, want: big.NewInt(0x2)},
		{input: `"0x2F2"`, want: big.NewInt(0x2f2)},
		{input: `"0X2F2"`, want: big.NewInt(0x2f2)},
		{input: `"0x1122aaff"`, want: big.NewInt(0x1122aaff)},
		{input: `"0xbBb"`, want: big.NewInt(0xbbb)},
		{input: `"0xfffffffff"`, want: big.NewInt(0xfffffffff)},
		{
			input: `"0x112233445566778899aabbccddeeff"`,
			want:  referenceBig("112233445566778899aabbccddeeff"),
		},
		{
			input: `"0xffffffffffffffffffffffffffffffffffff"`,
			want:  referenceBig("ffffffffffffffffffffffffffffffffffff"),
		},
		{
			input: `"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`,
			want:  referenceBig("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
	}

	// ------------------------
	// errors
	errJSONEOF = errors.New("unexpected end of JSON input")
)

// --------------------------
// package Test* json encoding
func TestUnmarshalBig(t *testing.T) {
	for _, test := range unmarshalBigTests {
		var v Big
		err := json.Unmarshal([]byte(test.input), &v)
		if !inl.CheckError(t, test.input, err, test.wantErr) {
			continue
		}
		if test.want != nil && test.want.(*big.Int).Cmp((*big.Int)(&v)) != 0 {
			t.Errorf("input %s: value mismatch: got %x, want %x", test.input, (*big.Int)(&v), test.want)
			continue
		}
	}
}

func BenchmarkUnmarshalBig(b *testing.B) {
	input := []byte(`"0x123456789abcdef123456789abcdef"`)
	for i := 0; i < b.N; i++ {
		var v Big
		if err := v.UnmarshalJSON(input); err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalBig(t *testing.T) {
	for _, test := range encodeBigTests {
		in := test.input.(*big.Int)
		out, err := json.Marshal((*Big)(in))
		if err != nil {
			t.Errorf("%d: %v", in, err)
			continue
		}
		if want := `"` + test.want + `"`; string(out) != want {
			t.Errorf("%d: MarshalJSON output mismatch: got %q, want %q", in, out, want)
			continue
		}
		if out := (*Big)(in).String(); out != test.want {
			t.Errorf("%x: String mismatch: got %q, want %q", in, out, test.want)
			continue
		}
	}
}

// ---------------------------
// package inner

func referenceBig(s string) *big.Int {
	b, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid")
	}
	return b
}

func referenceBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
