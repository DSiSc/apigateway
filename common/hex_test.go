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
	"bytes"
	"testing"

	inl "github.com/DSiSc/apigateway/common/internal"
	"github.com/pkg/errors"
)

// -------------------------------------
// Vars

// --------------------------------------
// Types

// ----------------------------------
// Function Test*

func TestEncodeToString(t *testing.T) {
	var encodeBytesTests = []struct {
		input interface{}
		want  string
	}{
		{[]byte{}, "0x"},
		{[]byte{0}, "0x00"},
		{[]byte{0, 0, 1, 2}, "0x00000102"},
	}
	for _, test := range encodeBytesTests {
		enc := Ghex.EncodeToString(test.input.([]byte))
		if enc != test.want {
			t.Errorf("input %x: wrong encoding %s", test.input, enc)
		}
	}
}

func TestDecodeString(t *testing.T) {
	var decodeBytesTests = []struct {
		input        string
		want         interface{}
		wantErr      error // if set, decoding must fail on any platform
		wantErr32bit error // if set, decoding must fail on 32bit platforms (used for Uint tests)
	}{

		// invalid
		{input: ``, wantErr: ErrEmptyData},
		{input: `0`, wantErr: ErrMissingPrefix},
		{input: `0x0`, wantErr: ErrOddLength},
		{input: `0x023`, wantErr: ErrOddLength},
		{input: `0xxx`, wantErr: errors.Errorf("invalid byte: U+0078 'x'")},
		{input: `0x0X`, wantErr: errors.Errorf("invalid byte: U+0058 'X'")},
		{input: `0x01zz01`, wantErr: errors.Errorf("invalid byte: U+007A 'z'")},
		// valid
		{input: `0x`, want: []byte{}},
		{input: `0X`, want: []byte{}},
		{input: `0x02`, want: []byte{0x02}},
		{input: `0X02`, want: []byte{0x02}},
		{input: `0xffffffffff`, want: []byte{0xff, 0xff, 0xff, 0xff, 0xff}},
		{
			input: `0xffffffffffffffffffffffffffffffffffff`,
			want:  []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
	}
	for _, test := range decodeBytesTests {
		dec, err := Ghex.DecodeString(test.input)
		if !inl.CheckError(t, test.input, err, test.wantErr) {
			continue
		}
		if !bytes.Equal(test.want.([]byte), dec) {
			t.Errorf("input %s: value mismatch: got %x, want %x", test.input, dec, test.want)
			continue
		}
	}

}

func TestDecode(t *testing.T) {

	decodeTests := []struct {
		input   []byte
		want    []byte
		wantErr error
	}{
		//valid
		{input: []byte(`0x0001`), want: []byte{0, 1}, wantErr: nil},
		{input: []byte(`0x0101`), want: []byte{1, 1}, wantErr: nil},
		//invalid
		{input: []byte(`00`), want: nil, wantErr: ErrMissingPrefix},
		{input: []byte(`0x0`), want: nil, wantErr: ErrOddLength},
		{input: []byte(``), want: nil, wantErr: ErrEmptyData},
		{input: []byte(`0xxx`), want: nil, wantErr: errors.Errorf("invalid byte: U+0078 'x'")},
	}

	for _, tcase := range decodeTests {

		dst := make([]byte, Ghex.DecodeLen(len(tcase.input)))
		_, err := Ghex.Decode(dst, tcase.input)
		if err != nil {
			if err.Error() != tcase.wantErr.Error() {
				t.Errorf("Error mismath: got %s, want %s", err.Error(), tcase.wantErr.Error())
			}
		} else {
			if !bytes.Equal(dst, tcase.want) {
				t.Errorf("Decode bytes mismath: got %v, want %v", dst, tcase.want)

			}
		}
	}
}

func TestEncode(t *testing.T) {

	decodeTests := []struct {
		input   []byte
		want    []byte
		wantErr error
		length  int
	}{
		//valid
		{input: []byte{0, 1}, want: []byte(`0x0001`), wantErr: nil, length: 6},
		{input: []byte{1, 1}, want: []byte(`0x0101`), wantErr: nil, length: 6},
		//invalid
		{input: []byte{1, 1}, want: []byte(`0x0101`), wantErr: errors.Errorf("encode dst length with 4, want 6"), length: 4},
		{input: []byte{1, 1}, want: []byte(`0x0101`), wantErr: errors.Errorf("encode dst length with 7, want 6"), length: 7},
	}

	for _, tcase := range decodeTests {

		dst := make([]byte, tcase.length)
		n, err := Ghex.Encode(dst, tcase.input)
		if err != nil && n == 0 {
			if err.Error() != tcase.wantErr.Error() {
				t.Errorf("Error mismath: got %s, want %s", err.Error(), tcase.wantErr.Error())
			}
		} else {
			if !bytes.Equal(dst, tcase.want) {
				t.Errorf("Encode bytes mismath: got %v, want %v", dst, tcase.want)
			}

		}
	}
}

// -------------------------------
// Function Bench*

func BenchmarkEncode(b *testing.B) {
	input := []byte{0, 0, 1, 2}
	for i := 0; i < b.N; i++ {
		Ghex.EncodeToString(input)
	}
}

// ---------------------------------
// Function inner
