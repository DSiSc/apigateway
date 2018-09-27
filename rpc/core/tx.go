package core

import (
	cmn "github.com/DSiSc/apigateway/common"
	"github.com/DSiSc/apigateway/core/types"
	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
	"github.com/DSiSc/blockchain"
	craft "github.com/DSiSc/craft/types"
	"github.com/DSiSc/txpool"
	wtypes "github.com/DSiSc/wallet/core/types"
	rpctypes "github.com/DSiSc/apigateway/rpc/core/types"
	"math/big"
	"errors"
)

var (
	swch chan<- interface{}
)

func SetSwCh(ch chan<- interface{}) {
	swch = ch
}

// ------------------------------
// package Consts, Vars

// SendTransaction will create a transaction from the given param, sign it and submit to the txpool.
//
// ```shell
// curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{see query exapmle}],"id":1}'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// args := '[{
//  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
//  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
//  "gas": "0x76c0",
//  "gasPrice": "0x9184e72a000",
//  "value": "0x9184e72a",
//  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
// }]'
// txHash, err := client.SendTransaction(params, true)
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
// 	"error": "",
// 	"result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
// 	"id": "1",
// 	"jsonrpc": "2.0"
// }
// ```
//
// Returns a transaction matching the given transaction hash.
//
// ### Query Parameters
//
// | Parameter | Type        | Default           | Required | Description                                                                                          |
// |-----------+-------------+-------------------+----------+------------------------------------------------------------------------------------------------------|
// | from      | DATA        | nil               | true     | The address the transaction is send from.                                                            |
// | to        | DATA        | nil               | Option   | The address the transaction is directed to.Option when creating new contract.                        |
// | gas       | QUANTITY    | 90000             | Option   | Integer of the gas provided for the transaction execution. It will return unused gas.                |
// | gasPrice  | QUANTITY    | To-Be-Determined  | true   | Integer of the gasPrice used for each paid gas.                                                      |
// | value     | QUANTITY    | nil               | Option   | Integer of the value sent with this transaction.                                                     |
// | data      | DATA        | nil               | Option     | The compiled code of a contract OR the hash of the invoked method signature and encoded parameters.  |
// | nonce     | QUANTITY    | nil               | Option   | Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.  |
//
// ### Query Example
//
// ```
// args: [{
//  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
//  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
//  "gas": "0x76c0", // 30400
//  "gasPrice": "0x9184e72a000", // 10000000000000
//  "value": "0x9184e72a", // 2441406250
//  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
// }]
// ```

// ### Returns
//
// - `hash`: `[]byte` - hash of the transaction
// - `error`: `error` - error detail info
func SendTransaction(args ctypes.SendTxArgs) (cmn.Hash, error) {
	// give an initValue when nonce is nil
	if args.Nonce == nil {
		args.Nonce = cmn.NewUint64(16)
	}
	// value can be nil
	var value *big.Int
	if args.Value != nil {
		value = args.Value.ToBigInt()
	} else {
		value = nil
	}
	// data can be nil
	var data []byte
	if args.Data != nil {
		data = args.Data.Bytes()
	} else {
		data = nil
	}
	// give an initValue when gas is nil
	var gas uint64
	if args.Gas != nil {
		gas = args.Gas.Touint64()
	} else {
		gas = uint64(0)
	}
	// new types.Transaction base on SendTxArgs
	tx := types.NewTransaction(
		args.Nonce.Touint64(),
		args.To,
		value,
		gas,
		args.GasPrice.ToBigInt(),
		data,
		types.BytesToAddress(args.From.Bytes()),
	)

	// SignTx
	key, _ := wtypes.DefaultTestKey()
	signer := new(wtypes.FrontierSigner)
	tx, err := wtypes.SignTx(tx, signer, key)
	if err != nil {
		return cmn.BytesToHash([]byte("Fail to signTx")), err
	}

	// Send Tx to gossipswith
	go func() {
		// send transacation to swch, wait for transaction ID
		var swMsg interface{}
		swMsg = tx
		swch <- swMsg
	}()

	txId := types.TxHash(tx)

	return (cmn.Hash)(txId), nil
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
/*type RPCTransaction struct {
	BlockHash        cmn.Hash       `json:"blockHash"`
	BlockNumber      cmn.Uint64      `json:"blockNumber"`
	From             *types.Address `json:"from"`
	Gas              cmn.Uint64   `json:"gas"`
	GasPrice         *cmn.Big       `json:"gasPrice"`
	Hash             *cmn.Hash       `json:"hash"`
	Input            cmn.Bytes      `json:"input"`
	Nonce            *cmn.Uint64     `json:"nonce"`
	To               *types.Address `json:"to"`
	TransactionIndex cmn.Uint       `json:"transactionIndex"`
	Value            *cmn.Big       `json:"value"`
	V                *cmn.Big       `json:"v"`
	R                *cmn.Big       `json:"r"`
	S                *cmn.Big       `json:"s"`
}*/

func GetTransactionByHash(hash cmn.Hash) (*rpctypes.RPCTransaction, error) {
	// Try to return an already finalized transaction
	bc, _ := blockchain.NewLatestStateBlockChain()
	if tx, blockHash, blockNumber, index, _ := bc.GetTransactionByHash(TypeConvert(&hash)); tx != nil {
		return newRPCTransaction(tx, (cmn.Hash)(blockHash), blockNumber, index)
	}
	// No finalized transaction, try to retrieve it from the pool
	txp := txpool.NewTxPool(txpool.DefaultTxPoolConfig)
	if tx := txp.GetTxByHash((craft.Hash)(hash)); tx != nil {
		return newRPCPendingTransaction(tx)
	}
	// Transaction unknown, return as such
	return nil, nil
}

func newRPCTransaction(tx *craft.Transaction, blockHash cmn.Hash, blockNumber uint64, index uint64) (*rpctypes.RPCTransaction, error) {
	var from *types.Address
	if tx.Data.From != nil {
		from = (*types.Address)(tx.Data.From)
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
		h := (cmn.Hash)(types.TxHash(tx))
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

	var to *types.Address
	if tx.Data.Recipient != nil {
		to = (*types.Address)(tx.Data.Recipient)
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

	result := &rpctypes.RPCTransaction{
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
	if blockHash != (cmn.Hash{}) {
		result.BlockHash = blockHash
		result.BlockNumber = cmn.Uint64(blockNumber)
		result.TransactionIndex = cmn.Uint(index)
	}
	return result, nil
}

// newRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func newRPCPendingTransaction(tx *craft.Transaction) (*rpctypes.RPCTransaction, error) {
	return newRPCTransaction(tx, cmn.Hash{}, 0, 0)
}

func GetTransactionReceipt(hash cmn.Hash) (*rpctypes.RPCReceipt, error) {
	// Try to return an already finalized transaction
	bc, _ := blockchain.NewLatestStateBlockChain()
	if tx, blockHash, blockNumber, index, _ := bc.GetTransactionByHash(TypeConvert(&hash)); tx != nil {
		if receipt, _, _, _, _ := bc.GetReceiptByTxHash(TypeConvert(&hash)); receipt != nil {
			return newRPCReceipt(tx, receipt, (cmn.Hash)(blockHash), blockNumber, index)
		}
	}
	// Receipt unknown, return as such
	return nil, nil
}

func newRPCReceipt(tx *craft.Transaction, receipt *craft.Receipt, blockHash cmn.Hash, blockNumber uint64, index uint64) (*rpctypes.RPCReceipt, error) {
	var hash *cmn.Hash
	if tx != nil {
		h := (cmn.Hash)(types.TxHash(tx))
		hash = &h
	} else {
		hash = nil
	}

	var from *types.Address
	if tx.Data.From != nil {
		from = (*types.Address)(tx.Data.From)
	} else {
		from = nil
	}

	var to *types.Address
	if tx.Data.Recipient != nil {
		to = (*types.Address)(tx.Data.Recipient)
	} else {
		to = nil
	}

	var root []byte
	if receipt.PostState != nil {
		root = receipt.PostState
	} else {
		root = nil
	}

	var status *cmn.Uint64
	if &receipt.Status != nil {
		s := cmn.Uint64(receipt.Status)
		status = &s
	} else {
		status = nil
	}

	var gasUsed *cmn.Uint64
	if &receipt.GasUsed != nil {
		g := cmn.Uint64(receipt.GasUsed)
		gasUsed = &g
	} else {
		gasUsed = nil
	}

	var cumulativeGasUsed *cmn.Uint64
	if &receipt.CumulativeGasUsed != nil {
		c := cmn.Uint64(receipt.CumulativeGasUsed)
		cumulativeGasUsed = &c
	} else {
		cumulativeGasUsed = nil
	}

	var logsBloom *craft.Bloom
	if &receipt.Bloom != nil {
		logsBloom = &receipt.Bloom
	} else {
		logsBloom = nil
	}

	var logs []*craft.Log
	if receipt.Logs != nil {
		logs = receipt.Logs
	} else {
		logs = nil
	}

	var contractAddress *types.Address
	if &receipt.ContractAddress != nil {
		contractAddress = (*types.Address)(&receipt.ContractAddress)
	} else {
		contractAddress = nil
	}

	result := &rpctypes.RPCReceipt{
		TransactionHash:   hash,
		From:              from,
		To:                to,
		Root:              root,
		Status:            status,
		GasUsed:           gasUsed,
		CumulativeGasUsed: cumulativeGasUsed,
		LogsBloom:         logsBloom,
		Logs:              logs,
		ContractAddress:   contractAddress,
	}
	if blockHash != (cmn.Hash{}) {
		result.BlockHash = blockHash
		result.BlockNumber = (*cmn.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = cmn.Uint(index)
	}
	return result, nil
}

func newRPCTransactionFromBlockIndex(b *craft.Block, index uint64) (*rpctypes.RPCTransaction, error) {
	txs := b.Transactions
	if index >= uint64(len(txs)) {
		return nil,errors.New("index is too large")
	}
	return newRPCTransaction(txs[index], (cmn.Hash)(b.HeaderHash), b.Header.Height, index)
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func GetTransactionByBlockHashAndIndex(blockHash cmn.Hash, index cmn.Uint) (*rpctypes.RPCTransaction, error) {
	bc, err := blockchain.NewLatestStateBlockChain()
	if block, _ := bc.GetBlockByHash(TypeConvert(&blockHash)); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil, err
}