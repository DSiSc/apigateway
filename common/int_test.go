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
)

// ----------------------------
// package Consts, Vars

var (
	encodeBigTests = []marshalTest{
		{referenceBig("0"), "0x0"},
		{referenceBig("1"), "0x1"},
		{referenceBig("ff"), "0xff"},
		{referenceBig("112233445566778899aabbccddeeff"), "0x112233445566778899aabbccddeeff"},
		{referenceBig("80a7f2c1bcc396c00"), "0x80a7f2c1bcc396c00"},
		{referenceBig("-80a7f2c1bcc396c00"), "-0x80a7f2c1bcc396c00"},
	}

	encodeUint64Tests = []marshalTest{
		{uint64(0), "0x0"},
		{uint64(1), "0x1"},
		{uint64(0xff), "0xff"},
		{uint64(0x1122334455667788), "0x1122334455667788"},
	}

	encodeUintTests = []marshalTest{
		{uint(0), "0x0"},
		{uint(1), "0x1"},
		{uint(0xff), "0xff"},
		{uint(0x11223344), "0x11223344"},
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

	errJSONEOF = errors.New("unexpected end of JSON input")
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

//func checkError(t *testing.T, input string, got, want error) bool {
//	if got == nil {
//		if want != nil {
//			t.Errorf("input %s: got no error, want %q", input, want)
//			return false
//		}
//		return true
//	}
//	if want == nil {
//		t.Errorf("input %s: unexpected error %q", input, got)
//	} else if got.Error() != want.Error() {
//		t.Errorf("input %s: got error %q, want %q", input, got, want)
//	}
//	return false
//}

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

func TestUnmarshalBig(t *testing.T) {
	for _, test := range unmarshalBigTests {
		var v Big
		err := json.Unmarshal([]byte(test.input), &v)
		if !checkError(t, test.input, err, test.wantErr) {
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

func TestToBigInt(t *testing.T) {

    input := referenceBig("0")
    want := NewBig(input).toBigInt()
    if input != want {
        t.Errorf("Error with to BigInt: got %d, want %d", input, want)
    }
    
}

var unmarshalUint64Tests = []unmarshalTest{
	// invalid encoding
	{input: "", wantErr: errJSONEOF},
	{input: "null", wantErr: errNonString(uint64T)},
	{input: "10", wantErr: errNonString(uint64T)},
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

func TestUnmarshalUint64(t *testing.T) {
	for _, test := range unmarshalUint64Tests {
		var v Uint64
		err := json.Unmarshal([]byte(test.input), &v)
		if !checkError(t, test.input, err, test.wantErr) {
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
		if out := (Uint64)(in).String(); out != test.want {
			t.Errorf("%x: String mismatch: got %q, want %q", in, out, test.want)
			continue
		}
	}
}

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
		if out := (Uint)(in).String(); out != test.want {
			t.Errorf("%x: String mismatch: got %q, want %q", in, out, test.want)
			continue
		}
	}
}

var (
	// These are variables (not constants) to avoid constant overflow
	// checks in the compiler on 32bit platforms.
	maxUint33bits = uint64(^uint32(0)) + 1
	maxUint64bits = ^uint64(0)
)

var unmarshalUintTests = []unmarshalTest{
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

func TestUnmarshalUint(t *testing.T) {
	for _, test := range unmarshalUintTests {
		var v Uint
		err := json.Unmarshal([]byte(test.input), &v)
		if uintBits == 32 && test.wantErr32bit != nil {
			checkError(t, test.input, err, test.wantErr32bit)
			continue
		}
		if !checkError(t, test.input, err, test.wantErr) {
			continue
		}
		if uint(v) != test.want.(uint) {
			t.Errorf("input %s: value mismatch: got %d, want %d", test.input, v, test.want)
			continue
		}
	}
}

