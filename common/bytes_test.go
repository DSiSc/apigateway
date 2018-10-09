package common

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -------------------------
// package Test*

// Test New
func TestNew(t *testing.T) {
	data := []byte("hello world")

	bz := NewBytes(data)

	assert.Equal(t, data, bz.Bytes())
}

// This is a trivial test for protobuf compatibility.
func TestMarshal(t *testing.T) {
	bz := []byte("hello world")
	dataB := Bytes(bz)
	bz2, err := dataB.Marshal()
	assert.Nil(t, err)
	assert.Equal(t, bz, bz2)

	var dataB2 Bytes
	err = (&dataB2).Unmarshal(bz)
	assert.Nil(t, err)
	assert.Equal(t, dataB, dataB2)
}

// Test that the hex encoding works.
func TestMarshalJSON(t *testing.T) {

	type TestStruct struct {
		B1 []byte
		B2 Bytes
	}

	cases := []struct {
		input    []byte
		expected string
	}{
		{[]byte(``), `{"B1":"","B2":"0x"}`},
		{[]byte(`a`), `{"B1":"YQ==","B2":"0x61"}`},
		{[]byte(`abc`), `{"B1":"YWJj","B2":"0x616263"}`},
		{[]byte{1}, `{"B1":"AQ==","B2":"0x01"}`},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			ts := TestStruct{B1: tc.input, B2: tc.input}

			// Test that it marshals correctly to JSON.
			jsonBytes, err := json.Marshal(ts)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.expected, string(jsonBytes))
		})
	}
}

// Test that the hex encoding works.
func TestUnmarshalJSON(t *testing.T) {

	type TestStruct struct {
		B1 []byte
		B2 Bytes
	}

	cases := []struct {
		input    string
		expected []byte
	}{
		{`{"B1":"","B2":""}`, []byte(``)},
		{`{"B1":"","B2":"0x"}`, []byte(``)},
		{`{"B1":"YQ==","B2":"0x61"}`, []byte(`a`)},
		{`{"B1":"YWJj","B2":"0x616263"}`, []byte(`abc`)},
		{`{"B1":"AQ==","B2":"0x01"}`, []byte{1}},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			ts := TestStruct{}
			err := json.Unmarshal([]byte(tc.input), &ts)
			if err != nil {
				t.Fatal(err)
			}
			//assert.Equal(t, ts.B1, tc.expected)
			assert.Equal(t, ts.B2, Bytes(tc.expected))
		})
	}
}
