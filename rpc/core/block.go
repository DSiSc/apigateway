package core

import (
	cmn "github.com/DSiSc/apigateway/common"
	apitypes "github.com/DSiSc/apigateway/core/types"
	rpctypes "github.com/DSiSc/apigateway/rpc/core/types"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/repository"
	"github.com/DSiSc/txpool"
)

//#### eth_getBlockByHash
//
//Returns information about a block by hash.
//
//
//##### Parameters
//
//1. `DATA`, 32 Bytes - Hash of a block.
//2. `Boolean` - If `true` it returns the full transaction objects, if `false` only the hashes of the transactions.
//
//```js
//params: [
//   '0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331',
//   true
//]
//```
//
//##### Returns
//
//`Object` - A block object, or `null` when no block was found:
//
//- `number`: `QUANTITY` - the block number. `null` when its pending block.
//- `hash`: `DATA`, 32 Bytes - hash of the block. `null` when its pending block.
//- `parentHash`: `DATA`, 32 Bytes - hash of the parent block.
//- `transactionsRoot`: `DATA`, 32 Bytes - the root of the transaction trie of the block.
//- `stateRoot`: `DATA`, 32 Bytes - the root of the final state trie of the block.
//- `receiptsRoot`: `DATA`, 32 Bytes - the root of the receipts trie of the block.
//- `miner`: `DATA`, 20 Bytes - the address of the beneficiary to whom the mining rewards were given.
//- `timestamp`: `QUANTITY` - the unix timestamp for when the block was collated.
//- `transactions`: `Array` - Array of transaction objects, or 32 Bytes transaction hashes depending on the last given parameter.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331", true],"id":1}'
//
//// Result
//{
//"id":1,
//"jsonrpc":"2.0",
//"result": {
//    "number": "0x1b4", // 436
//    "hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
//    "parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
//    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
//    "stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
//    "miner": "0x4e65fda2159562a496f9f3522f89122a3088497a",
//    "timestamp": "0x54e34e8e" // 1424182926
//    "transactions": [{...},{ ... }]
//  }
//}
//```
//
//***
func GetBlockByHash(blockHash cmn.Hash, fullTx bool) (*rpctypes.Blockdata, error) {
	bc, err := repository.NewLatestStateRepository()
	if err == nil {
		block, err := bc.GetBlockByHash(TypeConvert(&blockHash))
		if block != nil {
			return rpcOutputBlock(block, true, fullTx)
		}
		return nil, err
	}
	return nil, err
}

//#### eth_getBlockTransactionCountByHash
//
//Returns the number of transactions in a block from a block matching the given block hash.
//
//
//##### Parameters
//
//1. `DATA`, 32 Bytes - hash of a block.
//
//```js
//params: [
//   '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238'
//]
//```
//
//##### Returns
//
//`QUANTITY` - integer of the number of transactions in this block.
//
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByHash","params":["0xc94770007dda54cF92009BFF0dE90c06F603a09f"],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0xc" // 11
//}
//```
//
//***
func GetBlockTransactionCountByHash(blockHash cmn.Hash) (*cmn.Uint, error) {
	bc, err := repository.NewLatestStateRepository()
	if block, err := bc.GetBlockByHash(TypeConvert(&blockHash)); block != nil {
		n := cmn.Uint(len(block.Transactions))
		return &n, err
	}
	return nil, err
}

//#### eth_getBlockTransactionCountByNumber
//> >
//Returns the number of transactions in a block matching the given block number.
//
//
//##### Parameters
//
//1. `QUANTITY|TAG` - integer of a block number, or the string `"earliest"`, `"latest"` or `"pending"`, as in the [default block parameter](#the-default-block-parameter).
//
//```js
//params: [
//   '0xe8', // 232
//]
//```
//
//##### Returns
//
//`QUANTITY` - integer of the number of transactions in this block.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockTransactionCountByNumber","params":["0xe8"],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0xa" // 10
//}
//```
//
//***
func GetBlockTransactionCountByNumber(blockNr apitypes.BlockNumber) (*cmn.Uint, error) {
	bc, err := repository.NewLatestStateRepository()
	height := blockNr.Touint64()
	if block, err := bc.GetBlockByHeight(height); block != nil {
		n := cmn.Uint(len(block.Transactions))
		return &n, err
	}
	return nil, err
}

//#### eth_getBlockByNumber
//
//Returns information about a block by block number.
//
//##### Parameters
//
//1. `QUANTITY|TAG` - integer of a block number, or the string `"earliest"`, `"latest"` or `"pending"`, as in the [default block parameter](#the-default-block-parameter).
//2. `Boolean` - If `true` it returns the full transaction objects, if `false` only the hashes of the transactions.
//
//```js
//params: [
//   '0x1b4', // 436
//   true
//]
//```
//
//##### Returns
//
//See [eth_getBlockByHash](#eth_getblockbyhash)
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}'
//```
//
//Result see [eth_getBlockByHash](#eth_getblockbyhash)
//
//***
func GetBlockByNumber(blockNr apitypes.BlockNumber, fullTx bool) (*rpctypes.Blockdata, error) {
	bc, err := repository.NewLatestStateRepository()
	var block *types.Block
	if err == nil {
		if blockNr == apitypes.LatestBlockNumber {
			block = bc.GetCurrentBlock()
		} else {
			height := blockNr.Touint64()
			block, err = bc.GetBlockByHeight(height)
		}
		if block != nil {
			block.Header.GasLimit = 6721975
			return rpcOutputBlock(block, true, fullTx)
		}
		return nil, err
	}
	return nil, err
}

//#### eth_blockNumber
//
//Returns the number of most recent block.
//
//##### Parameters
//none
//
//##### Returns
//
//`QUANTITY` - integer of the current block number the client is on.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}'
//
//// Result
//{
//  "id":83,
//  "jsonrpc": "2.0",
//  "result": "0xc94" // 1207
//}
//```
//
//***
func BlockNumber() (*cmn.Uint64, error) {
	blockchain, err := repository.NewLatestStateRepository()
	if err == nil {
		blockHeight := blockchain.GetCurrentBlockHeight()
		lastHeight := (*cmn.Uint64)(&blockHeight)
		return lastHeight, err
	}
	return nil, err
}

//#### eth_getBalance
//
//Returns the balance of the account of given address.
//
//##### Parameters
//
//1. `DATA`, 20 Bytes - address to check for balance.
//2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](#the-default-block-parameter)
//
//```js
//params: [
//   '0xc94770007dda54cF92009BFF0dE90c06F603a09f',
//   'latest'
//]
//```
//
//##### Returns
//
//`QUANTITY` - integer of the current balance in wei.
//
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xc94770007dda54cF92009BFF0dE90c06F603a09f", "latest"],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0x0234c8a3397aab58" // 158972490234375000
//}
//```
//
//***
func GetBalance(address apitypes.Address, blockNr apitypes.BlockNumber) (*cmn.Big, error) {
	bc, err := repository.NewLatestStateRepository()
	var block *types.Block
	if err == nil {
		if blockNr == apitypes.LatestBlockNumber {
			block = bc.GetCurrentBlock()
		} else {
			height := blockNr.Touint64()
			block, err = bc.GetBlockByHeight(height)
		}
		if &block.HeaderHash != nil {
			bchash, errbc := repository.NewRepositoryByBlockHash(block.HeaderHash)
			if errbc == nil {
				balance := (bchash.GetBalance((types.Address)(address)))
				return (*cmn.Big)(balance), nil
			}
		}
		return nil, err
	}
	return nil, err
}

//#### eth_getCode
//
//Returns code at a given address.
//
//
//##### Parameters
//
//1. `DATA`, 20 Bytes - address.
//2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](#the-default-block-parameter).
//
//```js
//params: [
//   '0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b',
//   '0x2'  // 2
//]
//```
//
//##### Returns
//
//`DATA` - the code from the given address.
//
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getCode","params":["0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b", "0x2"],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0x600160008035811a818181146012578301005b601b6001356025565b8060005260206000f25b600060078202905091905056"
//}
//```
//
//***
func GetCode(address apitypes.Address, blockNr apitypes.BlockNumber) (*cmn.Bytes, error) {
	bc, err := repository.NewLatestStateRepository()
	var block *types.Block
	if err == nil {
		if blockNr == apitypes.LatestBlockNumber {
			block = bc.GetCurrentBlock()
		} else {
			height := blockNr.Touint64()
			block, err = bc.GetBlockByHeight(height)
		}
		if &block.HeaderHash != nil {
			bchash, errbc := repository.NewRepositoryByBlockHash(block.HeaderHash)
			if errbc == nil {
				code := (bchash.GetCode((types.Address)(address)))
				return cmn.NewBytes(code), nil
			}
		}
		return nil, err
	}
	return nil, err
}

//#### eth_getTransactionCount
//
//Returns the number of transactions *sent* from an address.
//
//
//##### Parameters
//
//1. `DATA`, 20 Bytes - address.
//2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](#the-default-block-parameter)
//
//```js
//params: [
//   '0xc94770007dda54cF92009BFF0dE90c06F603a09f',
//   'latest' // state at the latest block
//]
//```
//
//##### Returns
//
//`QUANTITY` - integer of the number of transactions send from this address.
//
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0xc94770007dda54cF92009BFF0dE90c06F603a09f","latest"],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0x1" // 1
//}
//```
//
//***
func GetTransactionCount(address apitypes.Address, blockNr apitypes.BlockNumber) (*cmn.Uint64, error) {

	if blockNr == apitypes.PendingBlockNumber {
		noncePool := txpool.GetPoolNonce((types.Address)(address))
		if noncePool > 0 {
			return (*cmn.Uint64)(&noncePool), nil
		}
	}
	bc, err := repository.NewLatestStateRepository()
	var block *types.Block
	if blockNr == apitypes.LatestBlockNumber || blockNr == apitypes.PendingBlockNumber {
		block = bc.GetCurrentBlock()
	} else {
		height := blockNr.Touint64()
		block, err = bc.GetBlockByHeight(height)
	}
	if &block.HeaderHash != nil {
		bchash, errbc := repository.NewRepositoryByBlockHash(block.HeaderHash)
		if errbc == nil {
			nonce := (bchash.GetNonce((types.Address)(address)))
			return (*cmn.Uint64)(&nonce), nil
		}
	}
	return nil, err
}

func TypeConvert(a *cmn.Hash) types.Hash {
	var hash types.Hash
	if a != nil {
		copy(hash[:], a[:])
	}
	return hash
}

func rpcOutputBlock(b *types.Block, inclTx bool, fullTx bool) (*rpctypes.Blockdata, error) {
	fields, err := RPCMarshalBlock(b, inclTx, fullTx)
	if err != nil {
		return nil, err
	}
	//fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(b.Hash()))
	return fields, err
}

func RPCMarshalBlock(b *types.Block, inclTx bool, fullTx bool) (*rpctypes.Blockdata, error) {
	head := b.Header // copies the header once
	fields := rpctypes.Blockdata{
		Number:           (cmn.Uint64)(head.Height),
		Hash:             (cmn.Hash)(b.HeaderHash),
		ParentHash:       (cmn.Hash)(head.PrevBlockHash),
		MixHash:          (cmn.Hash)(head.MixDigest),
		StateRoot:        (cmn.Hash)(head.StateRoot),
		Miner:            (apitypes.Address)(head.CoinBase),
		Timestamp:        (cmn.Uint64)(head.Timestamp),
		TransactionsRoot: (cmn.Hash)(head.TxRoot),
		ReceiptsRoot:     (cmn.Hash)(head.ReceiptsRoot),
		GasLimit:         (cmn.Uint64)(head.GasLimit),
	}

	if inclTx {
		var blockTxs []*rpctypes.BlockTransaction
		txs := b.Transactions
		for i := 0; i < len(txs); i++ {
			tx, _ := toRPCTransaction(txs[i])
			blockTxs = append(blockTxs, tx)
		}
		if fullTx {
			fields.Transactions = blockTxs
		}
	}

	return &fields, nil
}

func toRPCTransaction(tx *types.Transaction) (*rpctypes.BlockTransaction, error) {
	var from *apitypes.Address
	if tx.Data.From != nil {
		from = (*apitypes.Address)(tx.Data.From)
	} else {
		from = nil
	}

	var gas cmn.Uint64
	if &tx.Data.GasLimit != nil {
		gas = (cmn.Uint64)(tx.Data.GasLimit)
	} else {
		gas = cmn.Uint64(0)
	}

	var gasPrice *cmn.Big
	if tx.Data.Price != nil {
		gasPrice = (*cmn.Big)(tx.Data.Price)
	} else {
		gasPrice = nil
	}

	var hash *cmn.Hash
	if tx != nil {
		h := (cmn.Hash)(apitypes.TxHash(tx))
		hash = &h
	} else {
		hash = nil
	}

	var input cmn.Bytes
	if tx.Data.Payload != nil {
		input = cmn.Bytes(tx.Data.Payload)
	} else {
		input = nil
	}

	var nonce *cmn.Uint64
	if &tx.Data.AccountNonce != nil {
		nonce = (*cmn.Uint64)(&tx.Data.AccountNonce)
	} else {
		nonce = nil
	}

	var to *apitypes.Address
	if tx.Data.Recipient != nil {
		to = (*apitypes.Address)(tx.Data.Recipient)
	} else {
		to = nil
	}

	var value *cmn.Big
	if tx.Data.Amount != nil {
		value = (*cmn.Big)(tx.Data.Amount)
	} else {
		value = nil
	}

	var v *cmn.Big
	if tx.Data.V != nil {
		v = (*cmn.Big)(tx.Data.V)
	} else {
		v = nil
	}

	var r *cmn.Big
	if tx.Data.R != nil {
		r = (*cmn.Big)(tx.Data.R)
	} else {
		r = nil
	}

	var s *cmn.Big
	if tx.Data.S != nil {
		s = (*cmn.Big)(tx.Data.S)
	} else {
		s = nil
	}

	result := &rpctypes.BlockTransaction{
		From:     from,
		Gas:      gas,
		GasPrice: gasPrice,
		Hash:     hash,
		Input:    input,
		Nonce:    nonce,
		To:       to,
		Value:    value,
		V:        v,
		R:        r,
		S:        s,
	}

	return result, nil
}
