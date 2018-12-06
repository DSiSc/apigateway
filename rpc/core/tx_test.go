package core

import (
	"encoding/json"
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/txpool"
	"github.com/DSiSc/validator/worker"
	"github.com/DSiSc/validator/worker/common"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	cmn "github.com/DSiSc/apigateway/common"
	ctypes "github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/apigateway/rpc/lib/types"
	crafttypes "github.com/DSiSc/craft/types"
	wtypes "github.com/DSiSc/wallet/core/types"
	"github.com/stretchr/testify/assert"
)

// ------------------------
// package Consts, Vars

var hashtest = cmn.HexToHash("0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99")

// ------------------------
// package Test*

func TestSendTransaction(t *testing.T) {

	// -------------------------
	// Mock:  mockTransaction
	mockTransaction := ctypes.NewTransaction(uint64(11), &to, value, gas, gasPrice, data, from)
	mockContract := ctypes.NewTransaction(uint64(11), nil, nil, math.MaxUint64/2, new(big.Int).SetUint64(1), nil, from)

	// SignTx
	key, _ := wtypes.DefaultTestKey()
	mockTransaction, _ = wtypes.SignTx(mockTransaction, new(wtypes.FrontierSigner), key)
	mockContract, _ = wtypes.SignTx(mockContract, new(wtypes.FrontierSigner), key)

	// NOTE(peerlink): tx.hash changed when call tx.Hash()
	txId := ctypes.TxHash(mockTransaction)
	mockTxHash := ctypes.HashBytes(txId)

	ContractId := ctypes.TxHash(mockContract)
	mockContractHash := ctypes.HashBytes(ContractId)

	// -------------------------
	// set mock swch, before node start http server.
	mockSwCh := make(chan interface{})
	defer close(mockSwCh)
	SetSwCh(mockSwCh)

	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetNonce", func(*blockchain.BlockChain, crafttypes.Address) uint64 {
		return uint64(10)
	})

	monkey.Patch(txpool.GetPoolNonce, func(crafttypes.Address) uint64 {
		return uint64(10)
	})

	// ---------------------------
	// tests case
	tests := []struct {
		payload         string
		wantErr         string
		wantReturn      []byte
		mockTransaction *crafttypes.Transaction
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
			"", mockTxHash, mockTransaction},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendTransaction", "id": "0", "params": [{
              "from": "%s"}]}`, request.from),
			"", mockContractHash, mockContract},
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

		// --------------------
		// Test read from swch for SwitchMsg
		var swMsg interface{}
		swMsg = <-mockSwCh

		// assert: Type, Content.
		assert.Equal(t, reflect.TypeOf(tt.mockTransaction), reflect.TypeOf(swMsg), "swMsg type should be types.Transaction")

		// exceptData := new(big.Int).SetBytes(cmn.HexDecode("0x9184e72a000"))
		actualTransaction := swMsg.(*crafttypes.Transaction)
		assert.Equal(t, tt.mockTransaction, actualTransaction, "transaction price should equal request input")

		// ----------------
		// Test reponse
		recv := new(rpctypes.RPCResponse)
		json.Unmarshal(blob, recv)

		if tt.wantErr == "" {
			assert.Nil(t, recv.Error, "#%d: not expecting an error", i)
			// FIXME(peerlink): check return Hash and mockTransaction.Hash()
			var result cmn.Hash
			json.Unmarshal(recv.Result, &result)
			assert.Equal(t, tt.wantReturn, result.Bytes(), "Hash should equals")
		} else {
			assert.True(t, recv.Error.Code < 0, "#%d: not expecting a positive JSONRPC code", i)
			// The wanted error is either in the message or the data
			assert.Contains(t, recv.Error.Message+recv.Error.Data, tt.wantErr, "#%d: expected substring", i)
		}
	}

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetNonce")
	monkey.Unpatch(txpool.GetPoolNonce)
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetTransactionByHash(t *testing.T) {
	// ------------------------
	// mock
	mockReturnTx := getMockTx()
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash", func(*blockchain.BlockChain, crafttypes.Hash) (*crafttypes.Transaction, crafttypes.Hash, uint64, uint64, error) {
		return mockReturnTx, (crafttypes.Hash)(hashtest), 5, 7, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByHash", "id": 1, "params": [
              "0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0x5","from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0x88da6692896bd38bf3145afef63088632f74ca0c74b0221b8940ba8eb996d1f0","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x7","value":"0x0","v":"0x0","r":"0x0","s":"0x0"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)

}

func TestGetTransactionReceipt(t *testing.T) {
	// ------------------------
	// mock
	mockReturnTx := getMockTx()
	mockReturnReceipt := &crafttypes.Receipt{
		PostState:         getBytes("0x7b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e929"),
		Status:            uint64(2),
		CumulativeGasUsed: uint64(1239),
		Bloom:             crafttypes.Bloom{},
		Logs:              nil,
		TxHash:            (crafttypes.Hash)(hashtest),
		ContractAddress:   (crafttypes.Address)(ctypes.BytesToAddress(getBytes(request.from))),
		GasUsed:           uint64(1510),
	}
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash", func(*blockchain.BlockChain, crafttypes.Hash) (*crafttypes.Transaction, crafttypes.Hash, uint64, uint64, error) {
		return mockReturnTx, (crafttypes.Hash)(hashtest), 5, 7, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetReceiptByTxHash", func(*blockchain.BlockChain, crafttypes.Hash) (*crafttypes.Receipt, crafttypes.Hash, uint64, uint64, error) {
		return mockReturnReceipt, (crafttypes.Hash)(hashtest), 5, 7, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionReceipt", "id": 1, "params": [
              "0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0x5","transactionHash":"0x88da6692896bd38bf3145afef63088632f74ca0c74b0221b8940ba8eb996d1f0","transactionIndex":"0x7","from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","root":"e0ogr1SPXLN0gVeOE/bpYcUensG5k214HBBhMjmz6Sk=","status":"0x2","gasUsed":"0x5e6","cumulativeGasUsed":"0x4d7","logsBloom":null,"logs":null,"contractAddress":"0xb60e8dd61c5d32be8058bb8eb970870f07233155"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetReceiptByTxHash")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetTransactionByBlockHashAndIndex(t *testing.T) {
	// ------------------------
	// mock
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash", func(*blockchain.BlockChain, crafttypes.Hash) (*crafttypes.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByBlockHashAndIndex", "id": 1, "params": [
              "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b", "0x0"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0xc","from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xec415a5415907387bf6f24ba16409370d0b457548b0b2d4bca05c5bd5263e507","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x0","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)

}

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {
	// ------------------------
	// mock
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*crafttypes.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*blockchain.BlockChain) *crafttypes.Block {
		blockdata := getMockBlock()
		return blockdata
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByBlockNumberAndIndex", "id": 1, "params": [
              "0x1b4", "0x0"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0xc","from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xec415a5415907387bf6f24ba16409370d0b457548b0b2d4bca05c5bd5263e507","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x0","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}}`},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByBlockNumberAndIndex", "id": 1, "params": [
              "latest", "0x0"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0xc","from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xec415a5415907387bf6f24ba16409370d0b457548b0b2d4bca05c5bd5263e507","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x0","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)

}

func TestCall(t *testing.T) {
	// -------------------------
	// Mock:  mockTransaction
	mockTx1 := ctypes.NewTransaction(uint64(0), &to, value, gas, gasPrice, data, from)
	mockTx2 := ctypes.NewTransaction(uint64(0), &to, nil, math.MaxUint64/2, new(big.Int).SetUint64(1), nil, ctypes.Address{})
	// tests case
	tests := []struct {
		payload         string
		wantErr         string
		wantReturn      string
		mockTransaction *crafttypes.Transaction
	}{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_call", "id": 1, "params": [{
              "from": "%s",
              "to": "%s",
              "gas": "%s",
              "gasPrice": "%s",
              "value": "%s",
			  "nonce": "%s",
              "data": "%s"}, "0x1b4"]}`, request.from, request.to, request.gas, request.gasPrice,
				request.value, request.nonce, request.data),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x38"}`, mockTx1},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_call", "id": 1, "params": [{
              "from": "%s",
              "to": "%s"}, "0x1b4"]}`, request.from, request.to),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x38"}`, mockTx2},
	}

	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*crafttypes.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*blockchain.BlockChain) *crafttypes.Block {
		blockdata := getMockBlock()
		return blockdata
	})
	monkey.Patch(blockchain.NewBlockChainByBlockHash, func(crafttypes.Hash) (*blockchain.BlockChain, error) {
		return b, nil
	})

	mux := testMux()
	for i, tt := range tests {

		monkey.Patch(evm.NewEVMContext, func(tx crafttypes.Transaction, header *crafttypes.Header, bc *blockchain.BlockChain, addr crafttypes.Address) evm.Context {
			//assert.Equal(t, tt.mockTransaction, &tx)
			return evm.Context{}
		})

		monkey.Patch(evm.NewEVM, func(evm.Context, *blockchain.BlockChain) *evm.EVM {
			return &evm.EVM{}
		})

		monkey.Patch(worker.ApplyTransaction, func(*evm.EVM, *crafttypes.Transaction, *common.GasPool) ([]byte, uint64, bool, error) {
			return getBytes("0x38"), uint64(0), true, nil
		})

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
		recv := new(rpctypes.RPCResponse)
		json.Unmarshal(blob, recv)

		b, _ := json.Marshal(recv)
		assert.Equal(t, tt.wantReturn, string(b))

		monkey.Unpatch(worker.ApplyTransaction)
		monkey.Unpatch(evm.NewEVM)
		monkey.Unpatch(evm.NewEVMContext)

	}
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.Unpatch(blockchain.NewBlockChainByBlockHash)
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func getMockTx() *crafttypes.Transaction {

	recipient := (crafttypes.Address)(to)
	craftfrom := (crafttypes.Address)(from)
	d := crafttypes.TxData{
		AccountNonce: nonce,
		Recipient:    &recipient,
		From:         &craftfrom,
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     gas,
		Price:        gasPrice,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	mockReturnTx := &crafttypes.Transaction{Data: d}
	return mockReturnTx
}

func TestGasPrice(t *testing.T) {

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_gasPrice", "id": 1, "params": []}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x1"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

}

func TestEstimateGas(t *testing.T) {

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_estimateGas", "id": 1, "params": [{
              "from": "%s",
              "to": "%s",
              "gas": "%s",
              "gasPrice": "%s",
              "value": "%s",
			  "nonce": "%s",
              "data": "%s"
              }]}`, request.from, request.to, request.gas, request.gasPrice,
				request.value, request.nonce, request.data),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x64"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

}

func getBytes(input string) []byte {
	bytes, _ := cmn.Ghex.DecodeString(input)
	return bytes
}
