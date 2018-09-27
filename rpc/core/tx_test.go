package core

import (
	"encoding/json"
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/monkey"
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
	"github.com/DSiSc/apigateway/rpc/lib/types"
	crafttypes "github.com/DSiSc/craft/types"
	wtypes "github.com/DSiSc/wallet/core/types"
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
	requestContract = requestParams{
		from:     "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
		gas:      "0xbb8",
		gasPrice: "0x9184e72a0000",
		nonce:    "0x1",
		data:     "0x608060405234801561001057600080fd5b506040516105d93803806105d983398101604052805160008054600160a060020a0319163317905501805161004c906001906020840190610053565b50506100ee565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061009457805160ff19168380011785556100c1565b828001600101855582156100c1579182015b828111156100c15782518255916020019190600101906100a6565b506100cd9291506100d1565b5090565b6100eb91905b808211156100cd57600081556001016100d7565b90565b6104dc806100fd6000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166341c0e1b5811461005b5780634ac0d66e14610072578063cfae3217146100cb575b600080fd5b34801561006757600080fd5b50610070610155565b005b34801561007e57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100709436949293602493928401919081908401838280828437509497506101929650505050505050565b3480156100d757600080fd5b506100e0610382565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561011a578181015183820152602001610102565b50505050905090810190601f1680156101475780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60005473ffffffffffffffffffffffffffffffffffffffff163314156101905760005473ffffffffffffffffffffffffffffffffffffffff16ff5b565b806040518082805190602001908083835b602083106101c25780518252601f1990920191602091820191016101a3565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405180910390206001604051808280546001816001161561010002031660029004801561024c5780601f1061022a57610100808354040283529182019161024c565b820191906000526020600020905b815481529060010190602001808311610238575b505060408051918290038220818352600180546002600019610100838516150201909116049284018390529094507f047dcd1aa8b77b0b943642129c767533eeacd700c7c1eab092b8ce05d2b2faf59350918691819060208201906060830190869080156102fb5780601f106102d0576101008083540402835291602001916102fb565b820191906000526020600020905b8154815290600101906020018083116102de57829003601f168201915b5050838103825284518152845160209182019186019080838360005b8381101561032f578181015183820152602001610317565b50505050905090810190601f16801561035c5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a3805161037e906001906020840190610418565b5050565b60018054604080516020601f6002600019610100878916150201909516949094049384018190048102820181019092528281526060939092909183018282801561040d5780601f106103e25761010080835404028352916020019161040d565b820191906000526020600020905b8154815290600101906020018083116103f057829003601f168201915b505050505090505b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061045957805160ff1916838001178555610486565b82800160010185558215610486579182015b8281111561048657825182559160200191906001019061046b565b50610492929150610496565b5090565b61041591905b80821115610492576000815560010161049c5600a165627a7a723058202360802f45f120f2cd8bf9b7963e38317b72e805b49afc944fdce06a24372fd10029",
	}
	hashtest = cmn.HexToHash("0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99")
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

	mockTransaction := ctypes.NewTransaction(nonce, &to, value, gas, gasPrice, data, from)
	mockContract := ctypes.NewTransaction(nonce, nil, nil, 0, gasPrice, nil, from)


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

	// ---------------------------
	// tests case
	tests := []struct {
		payload string
		wantErr string
		wantReturn []byte
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
			"", mockTxHash, mockTransaction },
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_sendTransaction", "id": "0", "params": [{
              "from": "%s",
              "gasPrice": "%s"
              }]}`, request.from,  request.gasPrice),
			"", mockContractHash, mockContract },
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
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0x5","from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0x98a285ef7fb624025fd54d605455d0626b277eaf74206f2b31eb74949999788e","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionIndex":"0x7","value":"0x0","v":"0x0","r":"0x0","s":"0x0"}}`},
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
			"", `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","blockNumber":"0x5","transactionHash":"0x98a285ef7fb624025fd54d605455d0626b277eaf74206f2b31eb74949999788e","transactionIndex":"0x7","from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","root":"e0ogr1SPXLN0gVeOE/bpYcUensG5k214HBBhMjmz6Sk=","status":"0x2","gasUsed":"0x5e6","cumulativeGasUsed":"0x4d7","logsBloom":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==","logs":null,"contractAddress":"0xb60e8dd61c5d32be8058bb8eb970870f07233155"}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetReceiptByTxHash")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetTransactionByHash")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func getMockTx() *crafttypes.Transaction {
	nonce, _ := strconv.ParseUint(request.nonce[2:], 16, 32)
	to := ctypes.BytesToAddress(getBytes(request.to))
	from := ctypes.BytesToAddress(getBytes(request.from))
	gas, _ := strconv.ParseUint(request.gas[2:], 16, 32)
	//value := new(big.Int).SetBytes(getBytes(request.value))
	gasPrice := new(big.Int).SetBytes(getBytes(request.gasPrice))
	data := getBytes(request.data)
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

func getBytes(input string) []byte {
	bytes, _ := cmn.Ghex.DecodeString(input)
	return bytes
}
