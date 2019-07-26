package core

import (
	"encoding/json"
	"fmt"
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

	"encoding/hex"
	cmn "github.com/DSiSc/apigateway/common"
	ctypes "github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/DSiSc/craft/rlp"
	"github.com/DSiSc/craft/types"
	crafttypes "github.com/DSiSc/craft/types"
	"github.com/DSiSc/justitia/config"
	"github.com/DSiSc/repository"
	wtypes "github.com/DSiSc/wallet/core/types"
	"github.com/stretchr/testify/assert"
)

// ------------------------
// package Consts, Vars

var hashtest = cmn.HexToHash("0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99")

func TestETransaction_GetTxData(t *testing.T) {
	addr1 := ctypes.HexToAddress("0x59b3f85ba6eb737fd0fad93bc4b5f92fd8c591de")
	addr := ctypes.TypeConvert(&addr1)
	tx := &types.Transaction{
		Data: types.TxData{
			AccountNonce: 10,
			Price:        big.NewInt(1),
			GasLimit:     21000,
			Recipient:    addr,
			Amount:       big.NewInt(16),
		},
	}
	key, _ := wtypes.DefaultTestKey()
	signer := wtypes.NewEIP155Signer(big.NewInt(100))
	tx, _ = wtypes.SignTx(tx, signer, key)
	fmt.Printf("%x\n", ctypes.TxHash(tx))
}

// ------------------------
// package Test*

func TestSendTransaction(t *testing.T) {
	// -------------------------
	// Mock:  mockTransaction
	nonce := uint64(16)
	from := ctypes.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	to := ctypes.HexToAddress("0x59b3f85ba6eb737fd0fad93bc4b5f92fd8c591de")
	value := big.NewInt(0x10)
	gasLimit := uint64(0x5208)
	gasPrice := big.NewInt(0x1)
	data := "6080604052348015600f57600080fd5b50603580601d6000396000f3006080604052600080fd00a165627a7a72305820799fbb709d76627f42ecb481ef39ca4cebc2feb3c2314a3a9e3c2b872827cc7b0029"
	dataB, _ := hex.DecodeString(data)
	mockTransaction := ctypes.NewTransaction(nonce, &to, value, gasLimit, gasPrice, nil, from)
	mockContract := ctypes.NewTransaction(nonce+1, nil, nil, gasLimit, gasPrice, dataB, from)

	// SignTx
	key, _ := wtypes.DefaultTestKey()
	chainId, _ := config.GetChainIdFromConfig()
	signer := wtypes.NewEIP155Signer(big.NewInt(int64(chainId)))
	mockTransaction, _ = wtypes.SignTx(mockTransaction, signer, key)
	mockContract, _ = wtypes.SignTx(mockContract, signer, key)

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

	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetNonce", func(*repository.Repository, crafttypes.Address) uint64 {
		return mockTransaction.Data.AccountNonce + 1
	})

	monkey.Patch(txpool.GetPoolNonce, func(crafttypes.Address) uint64 {
		return mockTransaction.Data.AccountNonce + 1
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
              "from": "0x%x",
              "to": "0x%x",
              "gas": "0x%x",
              "gasPrice": "0x%x",
              "value": "0x%x",
			  "nonce": "0x%x"
              }]}`, from, to, gasLimit, gasPrice,
				value, nonce),
			"", mockTxHash, mockTransaction,
		},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendTransaction", "id": "0", "params": [{
             "from": "%s",
  			 "gas": "0x%x",
             "gasPrice": "0x%x",
			 "data":"0x%s"}]}`, request.from, gasLimit, gasPrice, data),
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
	monkey.Unpatch(repository.NewLatestStateRepository)
}

func TestSendRawTransaction(t *testing.T) {
	// -------------------------
	monkey.Patch(config.GetChainIdFromConfig, func() (uint64, error) { return 1, nil })
	defer monkey.UnpatchAll()
	// Mock:  mockTransaction
	nonce := uint64(16)
	from := ctypes.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	to := ctypes.HexToAddress("0x59b3f85ba6eb737fd0fad93bc4b5f92fd8c591de")
	value := big.NewInt(0x10)
	gasLimit := uint64(0x5208)
	gasPrice := big.NewInt(0x1)
	data := "6080604052348015600f57600080fd5b50603580601d6000396000f3006080604052600080fd00a165627a7a72305820799fbb709d76627f42ecb481ef39ca4cebc2feb3c2314a3a9e3c2b872827cc7b0029"
	dataB, _ := hex.DecodeString(data)
	mockTransaction := ctypes.NewTransaction(nonce, &to, value, gasLimit, gasPrice, make([]byte, 0), from)
	mockContract := ctypes.NewTransaction(nonce+1, nil, nil, gasLimit, gasPrice, dataB, from)

	// SignTx
	key, _ := wtypes.DefaultTestKey()
	chainId, _ := config.GetChainIdFromConfig()
	signer := wtypes.NewEIP155Signer(big.NewInt(int64(chainId)))
	mockTransaction, _ = wtypes.SignTx(mockTransaction, signer, key)
	wtypes.Sender(signer, mockTransaction)
	mockContract, _ = wtypes.SignTx(mockContract, signer, key)
	wtypes.Sender(signer, mockContract)

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

	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetNonce", func(*repository.Repository, crafttypes.Address) uint64 {
		return uint64(10)
	})

	monkey.Patch(txpool.GetPoolNonce, func(crafttypes.Address) uint64 {
		return uint64(10)
	})

	encodedMockTx, _ := rlp.EncodeToBytes(mockTransaction)
	mockTransaction.Size.Store(types.StorageSize(len(encodedMockTx)))
	encodedMockTxStr := fmt.Sprintf("0x%x", encodedMockTx)
	encodedMockContract, _ := rlp.EncodeToBytes(mockContract)
	mockContract.Size.Store(types.StorageSize(len(encodedMockContract)))
	encodedMockContractStr := fmt.Sprintf("0x%x", encodedMockContract)

	// ---------------------------
	// tests case
	tests := []struct {
		payload         string
		wantErr         string
		wantReturn      []byte
		mockTransaction *crafttypes.Transaction
	}{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendRawTransaction", "id": "0", "params": ["%s"]}`, encodedMockTxStr),
			"", mockTxHash, mockTransaction},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendRawTransaction", "id": "0", "params": ["%s"]}`, encodedMockContractStr),
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
	monkey.Unpatch(repository.NewLatestStateRepository)
}

func TestGetTransactionByHash(t *testing.T) {
	// ------------------------
	// mock
	mockReturnTx := getMockTx()
	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash", func(*repository.Repository, crafttypes.Hash) (*crafttypes.Transaction, crafttypes.Hash, uint64, uint64, error) {
		return mockReturnTx, (crafttypes.Hash)(hashtest), 5, 7, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByHash", "id": 1, "params": [
              "0x95d191e78062c420e863df03311e5a09b28b431ced6e65048362d65515cd5770"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0x5","from":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0x95d191e78062c420e863df03311e5a09b28b431ced6e65048362d65515cd5770","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x7","value":"0x0","v":"0x0","r":"0x0","s":"0x0"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash")
	monkey.Unpatch(repository.NewLatestStateRepository)

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
	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash", func(*repository.Repository, crafttypes.Hash) (*crafttypes.Transaction, crafttypes.Hash, uint64, uint64, error) {
		return mockReturnTx, (crafttypes.Hash)(hashtest), 5, 7, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetReceiptByTxHash", func(*repository.Repository, crafttypes.Hash) (*crafttypes.Receipt, crafttypes.Hash, uint64, uint64, error) {
		return mockReturnReceipt, (crafttypes.Hash)(hashtest), 5, 7, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionReceipt", "id": 1, "params": [
              "0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0x5","transactionHash":"0x95d191e78062c420e863df03311e5a09b28b431ced6e65048362d65515cd5770","transactionIndex":"0x7","from":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","root":"e0ogr1SPXLN0gVeOE/bpYcUensG5k214HBBhMjmz6Sk=","status":"0x2","gasUsed":"0x5e6","cumulativeGasUsed":"0x4d7","logsBloom":null,"logs":null,"contractAddress":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetReceiptByTxHash")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash")
	monkey.Unpatch(repository.NewLatestStateRepository)
}

func TestGetTransactionByBlockHashAndIndex(t *testing.T) {
	// ------------------------
	// mock
	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash", func(*repository.Repository, crafttypes.Hash) (*crafttypes.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByBlockHashAndIndex", "id": 1, "params": [
              "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b", "0x0"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0xc","from":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xbedd625a813484aca74b38242fd7f439735be6211a033bf088c8b7b3656f4192","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x0","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHash")
	monkey.Unpatch(repository.NewLatestStateRepository)

}

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {
	// ------------------------
	// mock
	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*repository.Repository, uint64) (*crafttypes.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*repository.Repository) *crafttypes.Block {
		blockdata := getMockBlock()
		return blockdata
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByBlockNumberAndIndex", "id": 1, "params": [
              "0x1b4", "0x0"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0xc","from":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xbedd625a813484aca74b38242fd7f439735be6211a033bf088c8b7b3656f4192","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x0","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}}`},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionByBlockNumberAndIndex", "id": 1, "params": [
              "latest", "0x0"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0xc","from":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xbedd625a813484aca74b38242fd7f439735be6211a033bf088c8b7b3656f4192","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x0","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock")
	monkey.Unpatch(repository.NewLatestStateRepository)

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

	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*repository.Repository, uint64) (*crafttypes.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*repository.Repository) *crafttypes.Block {
		blockdata := getMockBlock()
		return blockdata
	})
	monkey.Patch(repository.NewRepositoryByBlockHash, func(crafttypes.Hash) (*repository.Repository, error) {
		return b, nil
	})

	monkey.Patch(repository.NewLatestStateRepository, func() (*repository.Repository, error) {
		return b, nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetNonce", func(*repository.Repository, types.Address) uint64 {
		return uint64(0)
	})

	mux := testMux()
	for i, tt := range tests {

		monkey.Patch(evm.NewEVMContext, func(tx crafttypes.Transaction, header *crafttypes.Header, bc *repository.Repository, addr crafttypes.Address) evm.Context {
			//assert.Equal(t, tt.mockTransaction, &tx)
			return evm.Context{}
		})

		monkey.Patch(evm.NewEVM, func(evm.Context, *repository.Repository) *evm.EVM {
			return &evm.EVM{}
		})

		monkey.Patch(worker.ApplyTransaction, func(types.Address, *types.Header, *repository.Repository, *crafttypes.Transaction, *common.GasPool) ([]byte, uint64, bool, error, types.Address) {
			return getBytes("0x38"), uint64(0), true, nil, types.Address{}
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
	monkey.Unpatch(repository.NewRepositoryByBlockHash)
	monkey.Unpatch(repository.NewLatestStateRepository)
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
			"", `{"jsonrpc":"2.0","id":1,"result":"0x5208"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

}

func getBytes(input string) []byte {
	bytes, _ := cmn.Ghex.DecodeString(input)
	return bytes
}
