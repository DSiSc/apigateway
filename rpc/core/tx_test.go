package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	cmn "github.com/DSiSc/apigateway/common"
	ctypes "github.com/DSiSc/apigateway/core/types"
	rpctypes "github.com/DSiSc/apigateway/rpc/lib/types"
	crafttypes "github.com/DSiSc/craft/types"
	sw "github.com/DSiSc/gossipswitch/common"
	"github.com/stretchr/testify/assert"
)

// ------------------------
// package Consts, Vars

var (
	request = requestParams{
		from:     "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		to:       "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
		gas:      "0x76c0",
		gasPrice: "0x9184e72a0000",
		nonce:    "0x10",
		value:    "0x9184e72a",
		data:     "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
	}
)

// -------------------------
// package Struct

type requestParams struct {
	from     string
	to       string
	gas      string
	gasPrice string
	nonce    string
	value    string
	data     string
}

// ------------------------
// package Test*

func TestSendTransaction(t *testing.T) {

	// -------------------------
	// Mock:  mockTransaction
	nonce, _ := strconv.ParseUint(request.nonce[2:], 16, 32)
	to := ctypes.BytesToAddress(getBytes(request.to))
	from := ctypes.BytesToAddress(getBytes(request.from))
	gas, _ := strconv.ParseUint(request.gas[2:], 16, 32)
	value := new(big.Int).SetBytes(getBytes(request.value))
	gasPrice := new(big.Int).SetBytes(getBytes(request.gasPrice))
	data := getBytes(request.data)

	mockTransaction := ctypes.NewTransaction(nonce, to, value, gas, gasPrice, data, from)
	// NOTE(peerlink): tx.hash changed when call tx.Hash()
	ctypes.TxHash(mockTransaction)

	// -------------------------
	// set mock swch, before node start http server.
	mockSwCh := make(chan sw.SwitchMsg)
	defer close(mockSwCh)
	SetSwCh(mockSwCh)

	// ---------------------------
	// tests case
	tests := []struct {
		payload string
		wantErr string
	}{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendTransaction", "id": "0", "params": [{
              "from": "%s",
              "to": "%s",
              "gas": "%s",
              "gasPrice": "%s",
              "value": "%s",
			  "nonce": "%s",
              "data": "%s"
              }]}`, request.from, request.to, request.gas, request.gasPrice,
				request.value, request.nonce, request.data),
			""},
	}

	// ------------------------
	// httptest API
	mux := testMux()
	for i, tt := range tests {
		req, _ := http.NewRequest("POST", "http://localhost/", strings.NewReader(tt.payload))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// Test read from swch for SwitchMsg
		var swMsg sw.SwitchMsg
		swMsg = <-mockSwCh

		// assert: Type, Content.
		assert.Equal(t, reflect.TypeOf(mockTransaction), reflect.TypeOf(swMsg), "swMsg type should be types.Transaction")

		// exceptData := new(big.Int).SetBytes(cmn.HexDecode("0x9184e72a000"))
		actualTransaction := swMsg.(*crafttypes.Transaction)
		assert.Equal(t, mockTransaction, actualTransaction, "transaction price should equal request input")

		// --------------
		// Test Response
		res := rec.Result()
		// Always expecting back a JSONRPCResponse
		assert.True(t, statusOK(res.StatusCode), "#%d: should always return 2XX", i)
		blob, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("#%d: err reading body: %v", i, err)
			continue
		}

		recv := new(rpctypes.RPCResponse)
		json.Unmarshal(blob, recv)

		if tt.wantErr == "" {
			assert.Nil(t, recv.Error, "#%d: not expecting an error", i)
			// FIXME(peerlink): check return Hash and mockTransaction.Hash()
			//var result cmn.Hash
			//json.Unmarshal(recv.Result, &result)
			//assert.Equal(t, mockTxHash.Bytes(), result.Bytes(), "Hash should equals")
		} else {
			assert.True(t, recv.Error.Code < 0, "#%d: not expecting a positive JSONRPC code", i)
			// The wanted error is either in the message or the data
			assert.Contains(t, recv.Error.Message+recv.Error.Data, tt.wantErr, "#%d: expected substring", i)
		}
	}
}

func getBytes(input string) []byte {
	bytes, _ := cmn.Ghex.DecodeString(input)
	return bytes
}
