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

// Package hex utils functions.
package common

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/pkg/errors"
)

// ------------------------
// package Const, Vars
const (
	PREFIX = "0x"
)

var Ghex = New(PREFIX)

// errors
var (
	// E: Empty Data
	ErrEmptyData = errors.New("empty hex data")
	// E: without prefix
	ErrMissingPrefix = errors.Errorf("hex string without %s prefix", PREFIX)
	// E: odd length hex string
	ErrOddLength = errors.New("hex string of odd length")
)

// Hex
// Encode and Decode hex string with prefix 0x
type Hex struct {
	prefix []byte
}

func New(prefix string) *Hex {
	return &Hex{
		prefix: []byte(prefix),
	}
}

// DecodeLen decode dst len without prefix
func (h *Hex) DecodeLen(n int) int {
	if n <= len(h.prefix) {
		return hex.DecodedLen(len(h.prefix))
	}
	return hex.DecodedLen(n - len(h.prefix))
}

// Decode hex []byte to struct []byte
func (h *Hex) Decode(dst, src []byte) (int, error) {

	// verify data
	b, err := h.verify(src)
	if !b {
		return 0, err
	}

	// trim prefix
	nonePrefixdata := h.trimPrefix(src)

	// decode data to string
	n, err := hex.Decode(dst, nonePrefixdata)
	if err != nil {
		return n, wrapError(err)
	}

	return n, nil
}

// DecodeString decode hex string to []byte
func (h *Hex) DecodeString(data string) ([]byte, error) {

	bdata := []byte(data)

	// verify data
	b, err := h.verify(bdata)
	if !b {
		return bdata, err
	}

	// trim data prefix
	nonePrefixdata := h.trimPrefix(bdata)

	// decode data to string
	dstdata, err := hex.DecodeString(string(nonePrefixdata))
	if err != nil {
		return bdata, wrapError(err)
	}

	return dstdata, nil
}

func DecodeUint64(input string) (uint64, error) {
	raw, err := checkNumber(input)
	if err != nil {
		return 0, err
	}
	dec, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		return dec, wrapError(err)
	}
	return dec, err
}

func checkNumber(input string) (raw string, err error) {
	if len(input) == 0 {
		return "", ErrEmptyData
	}
	if !has0xPrefix(input) {
		return "", ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return "", ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return "", ErrLeadingZero
	}
	return input, nil
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// EncodeUint64 uint64 to []byte
func (h *Hex) EncodeUint64(data uint64) []byte {

	var buf = make([]byte, 2, 10)
	copy(buf, h.prefix)
	return strconv.AppendUint(buf, data, 16)
}

// EncodeUint uint to []byte
func (h *Hex) EncodeUint(data uint) []byte {
	return h.EncodeUint64(uint64(data))
}

// EncodeBig big.Int to []byte
func (h *Hex) EncodeBig(data *big.Int) []byte {
	nbits := data.BitLen()
	if nbits == 0 {
		return []byte(h.attachPrefix("0"))
	}

	return []byte(fmt.Sprintf("%#x", data))
}

// EncodeLen get dst []byte len
func (h *Hex) EncodeLen(n int) int {

	// EncodeLen + len(prefix)
	return hex.EncodedLen(n) + len(h.prefix)
}

func (h *Hex) Encode(dst, src []byte) (int, error) {

	if len(dst) == 0 || len(dst) != h.EncodeLen(len(src)) {
		return 0, errors.Errorf("encode dst length with %d, want %d", len(dst), h.EncodeLen(len(src)))
	}

	// encode [len(prefix):]
	n := hex.Encode(dst[len(h.prefix):], src)
	copy(dst, []byte(h.prefix))

	return n, nil
}

// EncodeToString encode []byte to string
func (h *Hex) EncodeToString(data []byte) string {
	dst := make([]byte, h.EncodeLen(len(data)))
	h.Encode(dst, data)
	return string(dst)
}

func (h *Hex) String() string {
	return fmt.Sprintf("Hex with prefix %s", string(h.prefix))
}

// MustDecode decode hex string to []byte or panic error
func (h *Hex) MustDecodeString(data string) []byte {
	bdata, err := h.DecodeString(data)
	if err != nil {
		panic(err)
	}
	return bdata
}

// HasPrefix check data has prefix
func (h *Hex) HasPrefix(data []byte) bool {
	return bytes.HasPrefix(data, h.prefix) || bytes.HasPrefix(data, bytes.ToUpper(h.prefix))
}

// TrimPrefix trim 0x prefix, data should with 0x prefix
func (h *Hex) trimPrefix(data []byte) []byte {
	return data[len(h.prefix):]
}

// attachPrefix attach prefix 0x
func (h *Hex) attachPrefix(data string) string {
	return string(h.prefix) + data
}

// Verify check data
func (h *Hex) verify(data []byte) (bool, error) {

	// "" ErrEmptyData
	if len(data) == 0 {
		return false, ErrEmptyData
	}

	// "0" ErrMissingPrefix
	if !h.HasPrefix(data) {
		return false, ErrMissingPrefix
	}

	return true, nil
}

// ----------------------
// package func inner
func wrapError(err error) error {

	wrapErr := err
	if err, ok := err.(hex.InvalidByteError); ok {
		wrapErr = errors.Errorf("invalid byte: %#U", rune(err))

	}
	if err == hex.ErrLength {
		wrapErr = ErrOddLength

	}

	return wrapErr
}
