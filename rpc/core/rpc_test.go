package core

import (
	"bytes"
	"encoding/json"
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	ctypes "github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/apigateway/log"
	rs "github.com/DSiSc/apigateway/rpc/lib/server"
	ltypes "github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/tendermint/go-amino"
)

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

	nonce, _ = strconv.ParseUint(request.nonce[2:], 16, 32)
	to       = ctypes.BytesToAddress(getBytes(request.to))
	from     = ctypes.BytesToAddress(getBytes(request.from))
	gas, _   = strconv.ParseUint(request.gas[2:], 16, 32)
	value    = new(big.Int).SetBytes(getBytes(request.value))
	gasPrice = new(big.Int).SetBytes(getBytes(request.gasPrice))
	data     = getBytes(request.data)
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

type Requestdata struct {
	payload    string
	wantErr    string
	wantRetrun string
}

func testMux() *http.ServeMux {

	cdc := amino.NewCodec()
	mux := http.NewServeMux()
	buf := new(bytes.Buffer)
	logger := log.NewTMLogger(buf)
	rs.RegisterRPCFuncs(mux, Routes, cdc, logger)

	return mux
}

func statusOK(code int) bool { return code >= 200 && code <= 299 }

func doRpcTest(t *testing.T, tests []*Requestdata) {
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

func getMockBlock() *types.Block {
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
		CoinBase:      (types.Address)(ctypes.BytesToAddress(getBytes("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b"))),
	}
	blockdata := types.Block{
		Header:       &headerdata,
		Transactions: transactions,
		HeaderHash:   (types.Hash)(hashtest),
	}
	return &blockdata
}
