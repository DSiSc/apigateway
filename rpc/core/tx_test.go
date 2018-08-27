package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	httpAddr = "http://127.0.0.1:47768"
)

func TestSendTransactionCurl(t *testing.T) {

	cases := []struct {
		input      string
		statusCode int
		want       interface{}
		wantErr    error
	}{
		// Test case 01: OK
		{
			input: `{"jsonrpc":"2.0","method":"eth_sendTransaction","params" : [{
  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
  "gas": "0x76c0",
  "gasPrice": "0x9184e72a000",
  "value": "0x9184e72a",
  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
}]
,"id":"1"}`,
			statusCode: 200,
			want:       "Not Implement",
			wantErr:    (*rpctypes.RPCError)(nil),
		},
		// Test case 02: Method not found
		{
			input: `{"jsonrpc":"2.0","method":"sendTransaction","params" : [{
  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
  "gas": "0x76c0",
  "gasPrice": "0x9184e72a000",
  "value": "0x9184e72a",
  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
}]
,"id":"1"}`,
			statusCode: 200,
			want:       "",
			wantErr:    &rpctypes.RPCError{Code: -32601, Message: "Method not found", Data: ""},
		},
		// Test case 03: from address not begin with "0x"
		{
			input: `{"jsonrpc":"2.0","method":"eth_sendTransaction","params" : [{
  "from": "01",
  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
  "gas": "0x76c0",
  "gasPrice": "0x9184e72a000",
  "value": "0x9184e72a",
  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
}]
,"id":"1"}`,
			statusCode: 200,
			want:       "",
			wantErr:    &rpctypes.RPCError{Code: -32602, Message: "Invalid params", Data: "Error converting json params to arguments: json: cannot unmarshal hex string without 0x prefix into Go value of type types.Address"},
		},
		// TODO(mengql): Test case 04: options
		//{
		//	input: `{"jsonrpc":"2.0","method":"eth_sendTransaction","params" : [{
		//  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
		//  "gas": "0x76c0",
		//  "gasPrice": "0x9184e72a000",
		//  "value": "0x9184e72a",
		//  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
		//}]
		//,"id":"1"}`,
		//	statusCode: 200,
		//	want:       "Not Implement",
		//	wantErr:    (*rpctypes.RPCError)(nil),
		//},
	}

	client := &http.Client{}

	for i, test := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			req, err := http.NewRequest("POST", httpAddr, strings.NewReader(test.input))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			require.Nil(t, err)

			resp, err := client.Do(req)
			require.Nil(t, err)

			defer resp.Body.Close()
			assert.Equal(t, test.statusCode, resp.StatusCode)

			responseBytes, _ := ioutil.ReadAll(resp.Body)
			response := &rpctypes.RPCResponse{}
			json.Unmarshal(responseBytes, response)

			assert.Equal(t, test.wantErr, response.Error)

			var result string
			json.Unmarshal(response.Result, &result)

			assert.Equal(t, test.want, result)
		})
	}
}
