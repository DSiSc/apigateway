package core

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/monkey"
	"math/big"
	"reflect"
	"testing"
)

// ------------------------

var b *blockchain.BlockChain

// -------------------------

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
			`{"jsonrpc":"2.0","id":1,"result":{"number":"0xc","hash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","parentHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","mixHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","stateRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","miner":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","timestamp":"0x85","transactionsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","receiptsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","transactions":[{"from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xec415a5415907387bf6f24ba16409370d0b457548b0b2d4bca05c5bd5263e507","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}]}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

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

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		blockdata := getMockBlock()
		return blockdata
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockByNumber", "id": 1, "params": [
              "0x1b4",true]}`),
			"",
			`{"jsonrpc":"2.0","id":1,"result":{"number":"0xc","hash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","parentHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","mixHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","stateRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","miner":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","timestamp":"0x85","transactionsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","receiptsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","transactions":[{"from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xec415a5415907387bf6f24ba16409370d0b457548b0b2d4bca05c5bd5263e507","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}]}}`},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBlockByNumber", "id": 1, "params": [
              "latest",true]}`),
			"",
			`{"jsonrpc":"2.0","id":1,"result":{"number":"0xc","hash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","parentHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","mixHash":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","stateRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","miner":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","timestamp":"0x85","transactionsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","receiptsRoot":"0x27b4a20af548f5cb37481578e13f6e961c51e9ec1b9936d781c10613239b3e99","transactions":[{"from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","gas":"0x76c0","gasPrice":"0x9184e72a0000","hash":"0xec415a5415907387bf6f24ba16409370d0b457548b0b2d4bca05c5bd5263e507","input":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675","nonce":"0x10","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","value":"0x9184e72a","v":"0x0","r":"0x0","s":"0x0"}]}}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock")
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
			"", `{"jsonrpc":"2.0","id":1,"result":"0x1"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

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
			"", `{"jsonrpc":"2.0","id":1,"result":"0x1"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestBlockNumber(t *testing.T) {
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlockHeight", func(*blockchain.BlockChain) uint64 {
		return uint64(56)
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_blockNumber", "id": 1, "params": []}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x38"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlockHeight")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetBalance(t *testing.T) {
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*types.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		blockdata := getMockBlock()
		return blockdata
	})
	monkey.Patch(blockchain.NewBlockChainByBlockHash, func(types.Hash) (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBalance", func(*blockchain.BlockChain, types.Address) *big.Int {
		return big.NewInt(int64(56))
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBalance", "id": 1, "params": ["0xc94770007dda54cF92009BFF0dE90c06F603a09f","latest"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x38"}`},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getBalance", "id": 1, "params": ["0xc94770007dda54cF92009BFF0dE90c06F603a09f","0x4"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x38"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBalance")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
	monkey.Unpatch(blockchain.NewBlockChainByBlockHash)
}

func TestGetCode(t *testing.T) {
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*types.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		blockdata := getMockBlock()
		return blockdata
	})
	monkey.Patch(blockchain.NewBlockChainByBlockHash, func(types.Hash) (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCode", func(*blockchain.BlockChain, types.Address) []byte {
		return []byte(`abc`)
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getCode", "id": 1, "params": ["0xc94770007dda54cF92009BFF0dE90c06F603a09f","latest"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x616263"}`},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getCode", "id": 1, "params": ["0xc94770007dda54cF92009BFF0dE90c06F603a09f","0x4"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x616263"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBalance")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
	monkey.Unpatch(blockchain.NewBlockChainByBlockHash)
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}

func TestGetTransactionCount(t *testing.T) {
	monkey.Patch(blockchain.NewLatestStateBlockChain, func() (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight", func(*blockchain.BlockChain, uint64) (*types.Block, error) {
		blockdata := getMockBlock()
		return blockdata, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		blockdata := getMockBlock()
		return blockdata
	})
	monkey.Patch(blockchain.NewBlockChainByBlockHash, func(types.Hash) (*blockchain.BlockChain, error) {
		return b, nil
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(b), "GetNonce", func(*blockchain.BlockChain, types.Address) uint64 {
		return uint64(57)
	})

	// tests case
	tests := []*Requestdata{
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionCount", "id": 1, "params": ["0xc94770007dda54cF92009BFF0dE90c06F603a09f","latest"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x39"}`},
		{

			fmt.Sprintf(`{"jsonrpc": "2.0", "method": "eth_getTransactionCount", "id": 1, "params": ["0xc94770007dda54cF92009BFF0dE90c06F603a09f","0x4"]}`),
			"", `{"jsonrpc":"2.0","id":1,"result":"0x39"}`},
	}
	// ------------------------
	// httptest API
	doRpcTest(t, tests)

	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetNonce")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetCurrentBlock")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(b), "GetBlockByHeight")
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
	monkey.Unpatch(blockchain.NewBlockChainByBlockHash)
	monkey.Unpatch(blockchain.NewLatestStateBlockChain)
}
