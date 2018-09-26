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
type Requestdata struct {
	payload string
	wantErr string
	wantRetrun string
}

var b *blockchain.BlockChain

// -------------------------
func getMockBlock() *types.Block{
	hashtest := cmn.HexToHash("0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99")
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
	return &blockdata
}

func DoRpcTest(t *testing.T, tests []*Requestdata) {
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
		assert.Equal(t, tt.wantRetrun, string(b))
	}
}

// ------------------------
// package Test*

func TestGetBlockByHash(t *testing.T) {
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash", func(*blockchain.BlockChain, types.Hash) (*types.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockByHash", "id": 1, "params": [
              "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d152733",true]}`),
			"",
			`{"jsonrpc":"2.0","id":1,"result":{"number":"0xc","hash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","parentHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","mixHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","stateRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","miner":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","timestamp":"0x85","transactionsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","receiptsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","transactions":[{"Data":{"nonce":"16","gasPrice":160000000000000,"gas":"30400","to":"1G6N1nxdMr6AWLuOuXCHDwckRWc=","from":"tg6N1hxdMr6AWLuOuXCHDwcjMVU=","value":2441406250,"input":"1G6N1nxdMr6NRujdZ8XTK+gFi7jrlwhw8HJEVnUFi7jrlwhw8HJEVnU=","v":0,"r":0,"s":0,"hash":null},"Hash":{},"Size":{},"From":{}}]}}`},
	}
	// ------------------------
	// httptest API
	DoRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetBlockByNumber(t *testing.T) {

	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*types.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	// tests case
	tests := []*Requestdata {
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockByNumber", "id": 1, "params": [
              "0x1b4",true]}`),
			"",
			`{"jsonrpc":"2.0","id":1,"result":{"number":"0xc","hash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","parentHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","mixHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","stateRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","miner":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","timestamp":"0x85","transactionsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","receiptsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","transactions":[{"Data":{"nonce":"16","gasPrice":160000000000000,"gas":"30400","to":"1G6N1nxdMr6AWLuOuXCHDwckRWc=","from":"tg6N1hxdMr6AWLuOuXCHDwcjMVU=","value":2441406250,"input":"1G6N1nxdMr6NRujdZ8XTK+gFi7jrlwhw8HJEVnUFi7jrlwhw8HJEVnU=","v":0,"r":0,"s":0,"hash":null},"Hash":{},"Size":{},"From":{}}]}}`},
	}
	// ------------------------
	// httptest API
	DoRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetBlockTransactionCountByHash(t *testing.T) {

	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash", func(*blockchain.BlockChain, types.Hash) (*types.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockTransactionCountByHash", "id": 1, "params": [
              "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d152733"]}`),
			"",`{"jsonrpc":"2.0","id":1,"result":"0x1"}`},
	}
	// ------------------------
	// httptest API
	DoRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetBlockTransactionCountByNumber(t *testing.T) {
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*types.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockTransactionCountByNumber", "id": 1, "params": [
              "0x1b4"]}`),
			"",`{"jsonrpc":"2.0","id":1,"result":"0x1"}`},
	}
	// ------------------------
	// httptest API
	DoRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}
