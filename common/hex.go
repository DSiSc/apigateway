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
	"encoding/hex"
	"strconv"
	"strings"
)

// ------------------------
// package Const, Vars
const (
	UINTBITS = 32 << (uint64(^uint(0)) >> 63)
	PREFIX   = "0x"
)

// Errors
var (
	ErrEmptyString   = NewError("empty hex string")
	ErrSyntax        = NewError("invalid hex string")
	ErrMissingPrefix = NewError("hex string without 0x prefix")
	ErrOddLength     = NewError("hex string of odd length")
	ErrEmptyNumber   = NewError("hex string \"0x\"")
	ErrLeadingZero   = NewError("hex number with leading zero digits")
)

// -----------------------
// package Functions

func HexEncode(srcBytes []byte) string {

	dstBytes := make([]byte, hex.EncodedLen(len(srcBytes)))

	hex.Encode(dstBytes, srcBytes)

	return PREFIX + string(dstBytes)
}

func HexDecode(srcStr string) ([]byte, error) {

	if len(srcStr) == 0 {
		return nil, ErrEmptyString
	}

	if !HexHasPrefix(srcStr) {
		return nil, ErrMissingPrefix
	}

	srcBytesWithoutPrefix := srcStr[len(PREFIX):]
	dstBytes, err := hex.DecodeString(srcBytesWithoutPrefix)
	if err != nil {
		err = explainError(err)
		return nil, err
	}

	return dstBytes, nil

}

func HexDecodeWithoutErr(srcStr string) []byte {

	dstBytes, _ := HexDecode(srcStr)
	return dstBytes
}

// Validate validates whether each byte is valid hexadecimal string.
func HexValidate(str string) bool {

	_, err := hex.DecodeString(str)
	if err != nil {
		return false
	}
	return true
}

func HexHasPrefix(str string) bool {
	// NOTE(peerlink): suport prefix: "0x", "0X"
	return strings.HasPrefix(str, PREFIX) || strings.HasPrefix(str, strings.ToUpper(PREFIX))
}

// -------------------
// pakcage Function inner

func explainError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return ErrUint64Range
		case strconv.ErrSyntax:
			return ErrSyntax
		}
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return ErrSyntax
	}
	if err == hex.ErrLength {
		return ErrOddLength
	}
	return err
}
