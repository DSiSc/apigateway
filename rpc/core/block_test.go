package core

import (
	"encoding/json"
	"fmt"
	cmn "github.com/DSiSc/apigateway/common"
	ctypes "github.com/DSiSc/apigateway/core/types"
	ltypes "github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// ------------------------
// package Consts, Vars

var ()

// -------------------------
// package Struct

// ------------------------
// package Test*

func TestGetBlockByHash(t *testing.T) {
	var b *blockchain.BlockChain
	hashtest := cmn.HexToHash("0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99")

	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash", func(*blockchain.BlockChain, types.Hash) (*types.Block, error) {

		nonce, _ := strconv.ParseUint(request.nonce[2:], 16, 32)
		to := ctypes.BytesToAddress(getBytes(request.to))
		from := ctypes.BytesToAddress(getBytes(request.from))
		gas, _ := strconv.ParseUint(request.gas[2:], 16, 32)
		value := new(big.Int).SetBytes(getBytes(request.value))
		gasPrice := new(big.Int).SetBytes(getBytes(request.gasPrice))
		data := getBytes(request.data)

		mockTransaction := ctypes.NewTransaction(nonce, &to, value, gas, gasPrice, data, from)
		transactions := make([]*types.Transaction, 1)
		transactions[0] = mockTransaction

		headerdata := types.Header{
			ChainID:       uint64(11),
			PrevBlockHash: (types.Hash)(hashtest),
			StateRoot:     (types.Hash)(hashtest),
			TxRoot:        (types.Hash)(hashtest),
			ReceiptsRoot:  (types.Hash)(hashtest),
			Height:        uint64(12),
			Timestamp:     uint64(133),
			MixDigest:     (types.Hash)(hashtest),
			Coinbase:      (types.Address)(ctypes.BytesToAddress(getBytes("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"))),
		}
		blockdata := types.Block{
			Header:       &headerdata,
			Transactions: transactions,
			HeaderHash:   (types.Hash)(hashtest),
		}
		return &blockdata, nil
	})

	block, err := GetBlockByHash(hashtest, true)
	var localback []byte
    //block is not nil
	assert.NotNil(t, block)
	assert.Nil(t, err)

	cdc := amino.NewCodec()
	a := ltypes.NewRPCSuccessResponse(cdc, 1, block)
	localback, _ = json.Marshal(a)

	// tests case
	tests := []struct {
		payload string
		wantErr string
	}{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockByHash", "id": 1, "params": [
              "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d152733",true]}`),
			""},
	}
	// ------------------------
	// httptest API
	mux := testMux()
	for i, tt := range tests {
		req, _ := http.NewRequest("POST", "http://localhost/", strings.NewReader(tt.payload))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

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
		// ----------------
		// Test reponse
		recv := new(ltypes.RPCResponse)
		json.Unmarshal(blob, recv)

		b, _ := json.Marshal(recv)

		if tt.wantErr == "" {
			assert.Nil(t, recv.Error, "#%d: not expecting an error", i)
			assert.Equal(t, string(b), string(localback), "Hash should equals")
		} else {
			assert.True(t, recv.Error.Code < 0, "#%d: not expecting a positive JSONRPC code", i)
			// The wanted error is either in the message or the data
			assert.Contains(t, recv.Error.Message+recv.Error.Data, tt.wantErr, "#%d: expected substring", i)
		}
	}

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetBlockByNumber(t *testing.T) {
	var b *blockchain.BlockChain
	hashtest := cmn.HexToHash("0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99")
	num := ctypes.BlockNumber(17)

	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*types.Block, error) {

		nonce, _ := strconv.ParseUint(request.nonce[2:], 16, 32)
		to := ctypes.BytesToAddress(getBytes(request.to))
		from := ctypes.BytesToAddress(getBytes(request.from))
		gas, _ := strconv.ParseUint(request.gas[2:], 16, 32)
		value := new(big.Int).SetBytes(getBytes(request.value))
		gasPrice := new(big.Int).SetBytes(getBytes(request.gasPrice))
		data := getBytes(request.data)

		mockTransaction := ctypes.NewTransaction(nonce, &to, value, gas, gasPrice, data, from)
		transactions := make([]*types.Transaction, 1)
		transactions[0] = mockTransaction

		headerdata := types.Header{
			ChainID:       uint64(11),
			PrevBlockHash: (types.Hash)(hashtest),
			StateRoot:     (types.Hash)(hashtest),
			TxRoot:        (types.Hash)(hashtest),
			ReceiptsRoot:  (types.Hash)(hashtest),
			Height:        uint64(12),
			Timestamp:     uint64(133),
			MixDigest:     (types.Hash)(hashtest),
			Coinbase:      (types.Address)(ctypes.BytesToAddress(getBytes("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"))),
		}
		blockdata := types.Block{
			Header:       &headerdata,
			Transactions: transactions,
			HeaderHash:   (types.Hash)(hashtest),
		}
		return &blockdata, nil
	})

	block, err := GetBlockByNumber(num, true)
	var localback []byte
	//block is not nil
    assert.NotNil(t, block)
	assert.Nil(t, err)

	cdc := amino.NewCodec()
	a := ltypes.NewRPCSuccessResponse(cdc, 1, block)
	localback, _ = json.Marshal(a)

	// tests case
	tests := []struct {
		payload string
		wantErr string
	}{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockByNumber", "id": 1, "params": [
              "0x1b4",true]}`),
			""},
	}
	// ------------------------
	// httptest API
	mux := testMux()
	for i, tt := range tests {
		req, _ := http.NewRequest("POST", "http://localhost/", strings.NewReader(tt.payload))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

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
		// ----------------
		// Test reponse
		recv := new(ltypes.RPCResponse)
		json.Unmarshal(blob, recv)

		b, _ := json.Marshal(recv)
		if tt.wantErr == "" {
			assert.Nil(t, recv.Error, "#%d: not expecting an error", i)
			assert.Equal(t, string(b), string(localback), "Hash should equals")
		} else {
			assert.True(t, recv.Error.Code < 0, "#%d: not expecting a positive JSONRPC code", i)
			// The wanted error is either in the message or the data
			assert.Contains(t, recv.Error.Message+recv.Error.Data, tt.wantErr, "#%d: expected substring", i)
		}
	}

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}