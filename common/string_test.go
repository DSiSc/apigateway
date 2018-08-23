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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("a", []string{"a", "b", "c"}))
	assert.False(t, StringInSlice("d", []string{"a", "b", "c"}))
	assert.True(t, StringInSlice("", []string{""}))
	assert.False(t, StringInSlice("", []string{}))
}

func TestIsHex(t *testing.T) {
	notHex := []string{
		"", "   ", "a", "x", "0", "0x", "0X", "0x ", "0X ", "0X a",
		"0xf ", "0x f", "0xp", "0x-",
		"0xf", "0XBED", "0xF", "0xbed", // Odd lengths
	}
	for _, v := range notHex {
		assert.False(t, IsHex(v), "%q is not hex", v)
	}
	hex := []string{
		"0x00", "0x0a", "0x0F", "0xFFFFFF", "0Xdeadbeef", "0x0BED",
		"0X12", "0X0A",
	}
	for _, v := range hex {
		assert.True(t, IsHex(v), "%q is hex", v)
	}
}

func TestSplitAndTrim(t *testing.T) {
	testCases := []struct {
		s        string
		sep      string
		cutset   string
		expected []string
	}{
		{"a,b,c", ",", " ", []string{"a", "b", "c"}},
		{" a , b , c ", ",", " ", []string{"a", "b", "c"}},
		{" a, b, c ", ",", " ", []string{"a", "b", "c"}},
		{" , ", ",", " ", []string{"", ""}},
		{"   ", ",", " ", []string{""}},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, SplitAndTrim(tc.s, tc.sep, tc.cutset), "%s", tc.s)
	}
}

func TestIsASCIIText(t *testing.T) {
	notASCIIText := []string{
		"", "\xC2", "\xC2\xA2", "\xFF", "\x80", "\xF0", "\n", "\t",
	}
	for _, v := range notASCIIText {
		assert.False(t, IsHex(v), "%q is not ascii-text", v)
	}
	asciiText := []string{
		" ", ".", "x", "$", "_", "abcdefg;", "-", "0x00", "0", "123",
	}
	for _, v := range asciiText {
		assert.True(t, IsASCIIText(v), "%q is ascii-text", v)
	}
}

func TestASCIITrim(t *testing.T) {
	assert.Equal(t, ASCIITrim(" "), "")
	assert.Equal(t, ASCIITrim(" a"), "a")
	assert.Equal(t, ASCIITrim("a "), "a")
	assert.Equal(t, ASCIITrim(" a "), "a")
	assert.Panics(t, func() { ASCIITrim("\xC2\xA2") })
}
