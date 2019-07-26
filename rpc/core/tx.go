package core

import (
	"errors"
	"fmt"
	acmn "github.com/DSiSc/apigateway/common"
	cmn "github.com/DSiSc/apigateway/common"
	"github.com/DSiSc/apigateway/core/types"
	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/monitor"
	"github.com/DSiSc/craft/rlp"
	craft "github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto"
	"github.com/DSiSc/justitia/config"
	"github.com/DSiSc/repository"
	"github.com/DSiSc/txpool"
	"github.com/DSiSc/validator/worker"
	"github.com/DSiSc/validator/worker/common"
	wtypes "github.com/DSiSc/wallet/core/types"
	"math"
	"math/big"
)

var (
	swch chan<- interface{}
)

func SetSwCh(ch chan<- interface{}) {
	swch = ch
}

//#### eth_sendTransaction
//
//Creates new message call transaction or a contract creation, if the data field contains code.
//
//##### Parameters
//
//1. `Object` - The transaction object
//- `from`: `DATA`, 20 Bytes - The address the transaction is send from.
//- `to`: `DATA`, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
//- `gas`: `QUANTITY`  - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
//- `gasPrice`: `QUANTITY`  - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas
//- `value`: `QUANTITY`  - (optional) Integer of the value sent with this transaction
//- `data`: `DATA`  - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters. For details see [Ethereum Contract ABI](https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI)
//- `nonce`: `QUANTITY`  - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
//
//```js
//params: [{
//  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
//  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
//  "gas": "0x76c0", // 30400
//  "gasPrice": "0x9184e72a000", // 10000000000000
//  "value": "0x9184e72a", // 2441406250
//  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
//}]
//```
//
//##### Returns
//
//`DATA`, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.
//
//Use [eth_getTransactionReceipt](#eth_gettransactionreceipt) to get the contract address, after the transaction was mined, when you created a contract.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{see above}],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331"
//}
//```
//
//***
func SendTransaction(args ctypes.SendTxArgs) (cmn.Hash, error) {
	monitor.JTMetrics.ApigatewayReceivedTx.Add(1)

	// give an initValue when nonce is nil
	var nonce uint64
	if args.Nonce == nil {
		bc, _ := repository.NewLatestStateRepository()
		noncePool := txpool.GetPoolNonce((craft.Address)(args.From))
		nonceChain := bc.GetNonce((craft.Address)(args.From))
		if noncePool > nonceChain {
			nonce = noncePool + 1
		} else {
			nonce = nonceChain
		}
	} else {
		nonce = args.Nonce.Touint64()
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
		if gas == 0 {
			gas = math.MaxUint64 / 2
		}
	} else {
		gas = math.MaxUint64 / 2
	}
	// give an initValue when gasPrice is nil
	var gasPrice *big.Int
	if args.GasPrice == nil {
		gasPrice = new(big.Int).SetUint64(types.DefaultGasPrice)
	} else {
		gasPrice = args.GasPrice.ToBigInt()
		if gasPrice.Sign() == 0 {
			gasPrice = new(big.Int).SetUint64(types.DefaultGasPrice)
		}
	}

	// new types.Transaction base on SendTxArgs
	tx := types.NewTransaction(
		nonce,
		args.To,
		value,
		gas,
		gasPrice,
		data,
		args.From,
	)

	// SignTx
	key, _ := wtypes.DefaultTestKey()
	chainId, err := config.GetChainIdFromConfig()
	if err != nil {
		log.Error("get chainId failed, err = ", err)
		return cmn.Hash{}, err
	}
	signer := wtypes.NewEIP155Signer(big.NewInt(int64(chainId)))
	tx, err = wtypes.SignTx(tx, signer, key)
	if err != nil {
		return cmn.BytesToHash([]byte("Fail to signTx")), err
	}
	txId := types.TxHash(tx)
	// Send Tx to gossip switch
	go func() {
		// send transaction to switch, wait for transaction ID
		var swMsg interface{}
		swMsg = tx
		swch <- swMsg
		monitor.JTMetrics.SwitchTakenTx.Add(1)
	}()
	return (cmn.Hash)(txId), nil
}

//#### eth_sendRawTransaction
//
//Creates new message call transaction or a contract creation for signed transactions.
//
//##### Parameters
//
//1. `DATA` - The signed transaction data.
//- `from`: `DATA`, 20 Bytes - The address the transaction is send from.
//- `to`: `DATA`, 20 Bytes - (optional when creating new contract) The address the transaction is directed to.
//- `gas`: `QUANTITY`  - (optional, default: 90000) Integer of the gas provided for the transaction execution. It will return unused gas.
//- `gasPrice`: `QUANTITY`  - (optional, default: To-Be-Determined) Integer of the gasPrice used for each paid gas
//- `value`: `QUANTITY`  - (optional) Integer of the value sent with this transaction
//- `data`: `DATA`  - The compiled code of a contract OR the hash of the invoked method signature and encoded parameters. For details see [Ethereum Contract ABI](https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI)
//- `nonce`: `QUANTITY`  - (optional) Integer of a nonce. This allows to overwrite your own pending transactions that use the same nonce.
//
//```js
//  params: ["0xf8acf8a70b869184e72a00008276c094d46e8dd67c5d32be8058bb8eb970870f0724456794b60e8dd61c5d32be8058bb8eb970870f07233155849184e72aa9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f0724456751ba0a81f79d00e342f5df1c47acabfd0ccc77a3f9ab919a15d5a6699d6de2c4ffbdda07b721e0eb6a50e7bce582dbd71f690004eab409abe7b6cb57b04a240d814ee6dc0c0c0"]
//```
//
//##### Returns
//
//`DATA`, 32 Bytes - the transaction hash, or the zero hash if the transaction is not yet available.

//
//Use [eth_getTransactionReceipt](#eth_gettransactionreceipt) to get the contract address, after the transaction was mined, when you created a contract.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":[{see above}],"id":1}'
//
//// Result
//{
//"id":1,
//"jsonrpc": "2.0",
//"result": "0x919d38fa5c395fa0f677e6554eef74fc7a48a64c087e320d538114c714d67d8f"
//}
//```
//
//***
func SendRawTransaction(encodedTx acmn.Bytes) (cmn.Hash, error) {
	monitor.JTMetrics.ApigatewayReceivedTx.Add(1)

	tx := new(craft.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		ethTx := new(craft.ETransaction)
		err = ethTx.DecodeBytes(encodedTx)
		if err != nil {
			log.Info("sendRawTransaction tx decode as ethereum error, err = ", err)
			return cmn.Hash{}, err
		}
		ethTx.SetTxData(&tx.Data)
	}

	//Caculate from and fill in Transaction
	chainId, err := config.GetChainIdFromConfig()
	if err != nil {
		log.Error("get chainId failed, err = ", err)
		return cmn.Hash{}, err
	}

	from, err := wtypes.Sender(wtypes.NewEIP155Signer(big.NewInt(int64(chainId))), tx)
	if err != nil {
		log.Error("get from address failed, err =  ", err)
		return cmn.Hash{}, err
	}
	from_tmp := craft.Address(from)
	tx.Data.From = &from_tmp
	// give an initValue when nonce is nil
	// Send Tx to gossip switch
	go func() {
		// send transaction to switch, wait for transaction ID
		var swMsg interface{}
		swMsg = tx
		swch <- swMsg
		monitor.JTMetrics.SwitchTakenTx.Add(1)
	}()
	txHash := types.TxHash(tx)
	return (cmn.Hash)(txHash), nil
}

func SendCrossRawTransaction(encodedTx acmn.Bytes, url string) (cmn.Hash, error) {
	monitor.JTMetrics.ApigatewayReceivedTx.Add(1)

	tx := new(craft.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return cmn.Hash{}, err
	}

	var args ctypes.SendTxArgs
	// need verify tx is legality提前验证一下

	// TODO：call the destination chain rpc(receiveCrossRawTransaction), like a sendTransaction
	receiveCrossRawTransactionReq(args)

	// return destination chain's tx hash
	return cmn.Hash{}, nil
}

func receiveCrossRawTransactionReq(args ctypes.SendTxArgs) (cmn.Hash, error) {
	monitor.JTMetrics.ApigatewayReceivedTx.Add(1)
	// Patchwork tx，fix from -- get publicAccount
	addr, err := getPubliceAcccount()
	if err != nil {
		return cmn.Hash{}, err
	}
	crossFrom := args.From

	args.From = addr
	// like sendTransaction, need sig
	private := "29ad43a4ebb4a65436d9fb116d471d96516b3d5cc153e045b384664bed5371b9"

	//get nonce
	bc, _ := repository.NewLatestStateRepository()
	noncePool := txpool.GetPoolNonce((craft.Address)(args.From))
	nonceChain := bc.GetNonce((craft.Address)(args.From))
	nonce := uint64(0)
	if noncePool > nonceChain {
		nonce = noncePool + 1
	} else {
		nonce = nonceChain
	}

	//amount & gasPrice
	amount := args.Value.ToBigInt()
	gasPrice := args.GasPrice.ToBigInt()
	tx := types.NewTransaction(uint64(nonce), args.To, amount, args.Gas.Touint64(), gasPrice, crossFrom.Bytes(), args.From)

	//sign tx
	priKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return cmn.Hash{}, err
	}
	chainID := big.NewInt(1)
	wtypes.SignTx(tx, wtypes.NewEIP155Signer(chainID), priKey)

	return cmn.Hash{}, nil
}

func getPubliceAcccount() (types.Address, error) {
	//get from config or genesis ?
	addr := "0x0fA3E9c7065Cf9b5f513Fb878284f902d167870c"
	address := types.HexToAddress(addr)

	return address, nil
}

//#### eth_getTransactionByHash
//
//Returns the information about a transaction requested by transaction hash.
//
//
//##### Parameters
//
//1. `DATA`, 32 Bytes - hash of a transaction
//
//```js
//params: [
//   "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"
//]
//```
//
//##### Returns
//
//`Object` - A transaction object, or `null` when no transaction was found:
//
//- `blockHash`: `DATA`, 32 Bytes - hash of the block where this transaction was in. `null` when its pending.
//- `blockNumber`: `QUANTITY` - block number where this transaction was in. `null` when its pending.
//- `from`: `DATA`, 20 Bytes - address of the sender.
//- `gas`: `QUANTITY` - gas provided by the sender.
//- `gasPrice`: `QUANTITY` - gas price provided by the sender in Wei.
//- `hash`: `DATA`, 32 Bytes - hash of the transaction.
//- `input`: `DATA` - the data send along with the transaction.
//- `nonce`: `QUANTITY` - the number of transactions made by the sender prior to this one.
//- `to`: `DATA`, 20 Bytes - address of the receiver. `null` when its a contract creation transaction.
//- `transactionIndex`: `QUANTITY` - integer of the transactions index position in the block. `null` when its pending.
//- `value`: `QUANTITY` - value transferred in Wei.
//- `v`: `QUANTITY` - ECDSA recovery id
//- `r`: `DATA`, 32 Bytes - ECDSA signature r
//- `s`: `DATA`, 32 Bytes - ECDSA signature s
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"],"id":1}'
//
//// Result
//{
//  "jsonrpc":"2.0",
//  "id":1,
//  "result":{
//    "blockHash":"0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
//    "blockNumber":"0x5daf3b", // 6139707
//    "from":"0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
//    "gas":"0xc350", // 50000
//    "gasPrice":"0x4a817c800", // 20000000000
//    "hash":"0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b",
//    "input":"0x68656c6c6f21",
//    "nonce":"0x15", // 21
//    "to":"0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
//    "transactionIndex":"0x41", // 65
//    "value":"0xf3dbb76162000", // 4290000000000000
//    "v":"0x25", // 37
//    "r":"0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea",
//    "s":"0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c"
//  }
//}
//```
//***
func GetTransactionByHash(hash cmn.Hash) (*ctypes.RPCTransaction, error) {
	// Try to return an already finalized transaction
	bc, _ := repository.NewLatestStateRepository()
	if tx, blockHash, blockNumber, index, _ := bc.GetTransactionByHash(TypeConvert(&hash)); tx != nil {
		return newRPCTransaction(tx, (cmn.Hash)(blockHash), blockNumber, index)
	}
	// No finalized transaction, try to retrieve it from the pool
	if tx := txpool.GetTxByHash((craft.Hash)(hash)); tx != nil {
		return newRPCPendingTransaction(tx)
	}
	// Transaction unknown, return as such
	return nil, nil
}

func newRPCTransaction(tx *craft.Transaction, blockHash cmn.Hash, blockNumber uint64, index uint64) (*ctypes.RPCTransaction, error) {
	var from *types.Address
	from = (*types.Address)(tx.Data.From)

	var gas cmn.Uint64
	gas = (cmn.Uint64)(tx.Data.GasLimit)

	var gasPrice *cmn.Big
	gasPrice = (*cmn.Big)(tx.Data.Price)

	var hash *cmn.Hash
	h := (cmn.Hash)(types.TxHash(tx))
	hash = &h

	var input cmn.Bytes
	input = cmn.Bytes(tx.Data.Payload)

	var nonce *cmn.Uint64
	nonce = (*cmn.Uint64)(&tx.Data.AccountNonce)

	var to *types.Address
	to = (*types.Address)(tx.Data.Recipient)

	var value *cmn.Big
	value = (*cmn.Big)(tx.Data.Amount)

	var v *cmn.Big
	v = (*cmn.Big)(tx.Data.V)

	var r *cmn.Big
	r = (*cmn.Big)(tx.Data.R)

	var s *cmn.Big
	s = (*cmn.Big)(tx.Data.S)

	result := &ctypes.RPCTransaction{
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

func newRPCPendingTransaction(tx *craft.Transaction) (*ctypes.RPCTransaction, error) {
	return newRPCTransaction(tx, cmn.Hash{}, 0, 0)
}

//
//#### eth_getTransactionReceipt
//
//Returns the receipt of a transaction by transaction hash.
//
//**Note** That the receipt is not available for pending transactions.
//
//
//##### Parameters
//
//1. `DATA`, 32 Bytes - hash of a transaction
//
//```js
//params: [
//   '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238'
//]
//```
//
//##### Returns
//
//`Object` - A transaction receipt object, or `null` when no receipt was found:
//
//- `transactionHash `: `DATA`, 32 Bytes - hash of the transaction.
//- `transactionIndex`: `QUANTITY` - integer of the transactions index position in the block.
//- `blockHash`: `DATA`, 32 Bytes - hash of the block where this transaction was in.
//- `blockNumber`: `QUANTITY` - block number where this transaction was in.
//- `from`: `DATA`, 20 Bytes - address of the sender.
//- `to`: `DATA`, 20 Bytes - address of the receiver. null when its a contract creation transaction.
//- `cumulativeGasUsed `: `QUANTITY ` - The total amount of gas used when this transaction was executed in the block.
//- `gasUsed `: `QUANTITY ` - The amount of gas used by this specific transaction alone.
//- `contractAddress `: `DATA`, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise `null`.
//- `logs`: `Array` - Array of log objects, which this transaction generated.
//- `logsBloom`: `DATA`, 256 Bytes - Bloom filter for light clients to quickly retrieve related logs.
//
//It also returns _either_ :
//
//- `root` : `DATA` 32 bytes of post-transaction stateroot (pre Byzantium)
//- `status`: `QUANTITY` either `1` (success) or `0` (failure)
//
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"],"id":1}'
//
//// Result
//{
//"id":1,
//"jsonrpc":"2.0",
//"result": {
//     transactionHash: '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238',
//     transactionIndex:  '0x1', // 1
//     blockNumber: '0xb', // 11
//     blockHash: '0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b',
//     cumulativeGasUsed: '0x33bc', // 13244
//     gasUsed: '0x4dc', // 1244
//     contractAddress: '0xb60e8dd61c5d32be8058bb8eb970870f07233155', // or null, if none was created
//     logs: [{
//         // logs as returned by getFilterLogs, etc.
//     }, ...],
//     logsBloom: "0x00...0", // 256 byte bloom filter
//     status: '0x1'
//  }
//}
//```
//***
func GetTransactionReceipt(hash cmn.Hash) (*ctypes.RPCReceipt, error) {
	// Try to return an already finalized transaction
	bc, _ := repository.NewLatestStateRepository()
	if tx, blockHash, blockNumber, index, _ := bc.GetTransactionByHash(TypeConvert(&hash)); tx != nil {
		if receipt, _, _, _, _ := bc.GetReceiptByTxHash(TypeConvert(&hash)); receipt != nil {
			return newRPCReceipt(tx, receipt, (cmn.Hash)(blockHash), blockNumber, index)
		}
	}
	// Receipt unknown, return as such
	return nil, nil
}

func newRPCReceipt(tx *craft.Transaction, receipt *craft.Receipt, blockHash cmn.Hash, blockNumber uint64, index uint64) (*ctypes.RPCReceipt, error) {
	var hash *cmn.Hash
	h := (cmn.Hash)(types.TxHash(tx))
	hash = &h

	var from *types.Address
	from = (*types.Address)(tx.Data.From)

	var to *types.Address
	to = (*types.Address)(tx.Data.Recipient)

	var root []byte
	if receipt.PostState != nil {
		root = receipt.PostState
	} else {
		root = nil
	}

	var status *cmn.Uint64
	s := cmn.Uint64(receipt.Status)
	status = &s

	var gasUsed *cmn.Uint64
	g := cmn.Uint64(receipt.GasUsed)
	gasUsed = &g

	var cumulativeGasUsed *cmn.Uint64
	c := cmn.Uint64(receipt.CumulativeGasUsed)
	cumulativeGasUsed = &c

	var logsBloom []byte
	copy(logsBloom[:], receipt.Bloom[:])

	var logs []*craft.Log
	if receipt.Logs != nil {
		logs = receipt.Logs
	} else {
		logs = nil
	}

	var contractAddress *types.Address
	if receipt.ContractAddress != (craft.Address{}) {
		contractAddress = (*types.Address)(&receipt.ContractAddress)
	} else {
		contractAddress = nil
	}

	result := &ctypes.RPCReceipt{
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

func newRPCTransactionFromBlockIndex(b *craft.Block, index uint64) (*ctypes.RPCTransaction, error) {
	txs := b.Transactions
	if index >= uint64(len(txs)) {
		return nil, errors.New("index is too large")
	}
	return newRPCTransaction(txs[index], (cmn.Hash)(b.HeaderHash), b.Header.Height, index)
}

//#### eth_getTransactionByBlockHashAndIndex
//
//Returns information about a transaction by block hash and transaction index position.
//
//
//##### Parameters
//
//1. `DATA`, 32 Bytes - hash of a block.
//2. `QUANTITY` - integer of the transaction index position.
//
//```js
//params: [
//   '0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331',
//   '0x0' // 0
//]
//```
//
//##### Returns
//
//See [eth_getTransactionByHash](#eth_gettransactionbyhash)
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByBlockHashAndIndex","params":["0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b", "0x0"],"id":1}'
//```

//Result see [eth_getTransactionByHash](#eth_gettransactionbyhash)
//***
func GetTransactionByBlockHashAndIndex(blockHash cmn.Hash, index cmn.Uint) (*ctypes.RPCTransaction, error) {
	bc, err := repository.NewLatestStateRepository()
	if block, _ := bc.GetBlockByHash(TypeConvert(&blockHash)); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil, err
}

//#### eth_getTransactionByBlockNumberAndIndex
//
//Returns information about a transaction by block number and transaction index position.
//
//
//##### Parameters
//
//1. `QUANTITY|TAG` - a block number, or the string `"earliest"`, `"latest"` or `"pending"`, as in the [default block parameter](#the-default-block-parameter).
//2. `QUANTITY` - the transaction index position.
//
//```js
//params: [
//   '0x29c', // 668
//   '0x0' // 0
//]
//```
//
//##### Returns
//
//See [eth_getTransactionByHash](#eth_gettransactionbyhash)
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params":["0x29c", "0x0"],"id":1}'
//```
//
//Result see [eth_getTransactionByHash](#eth_gettransactionbyhash)
//
//***
func GetTransactionByBlockNumberAndIndex(blockNr types.BlockNumber, index cmn.Uint) (*ctypes.RPCTransaction, error) {
	bc, err := repository.NewLatestStateRepository()
	var block *craft.Block
	if blockNr == types.LatestBlockNumber {
		block = bc.GetCurrentBlock()
	} else {
		height := blockNr.Touint64()
		block, err = bc.GetBlockByHeight(height)
	}
	if block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil, err
}

//#### eth_call
//
//Executes a new message call immediately without creating a transaction on the block chain.
//
//
//##### Parameters
//
//1. `Object` - The transaction call object
//- `from`: `DATA`, 20 Bytes - (optional) The address the transaction is sent from.
//- `to`: `DATA`, 20 Bytes  - The address the transaction is directed to.
//- `gas`: `QUANTITY`  - (optional) Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.
//- `gasPrice`: `QUANTITY`  - (optional) Integer of the gasPrice used for each paid gas
//- `value`: `QUANTITY`  - (optional) Integer of the value sent with this transaction
//- `data`: `DATA`  - (optional) Hash of the method signature and encoded parameters. For details see [Ethereum Contract ABI](https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI)
//2. `QUANTITY|TAG` - integer block number, or the string `"latest"`, `"earliest"` or `"pending"`, see the [default block parameter](#the-default-block-parameter)
//
//##### Returns
//
//`DATA` - the return value of executed contract.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_call","params":[{see above}],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0x"
//}
//```
//
//***
func Call(args ctypes.SendTxArgs, blockNr types.BlockNumber) (cmn.Bytes, error) {
	// to can not be nil
	var to *types.Address
	if &args.To == nil || *args.To == (types.Address{}) {
		return nil, errors.New("to is nil")
	} else {
		to = args.To
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
		if gas == 0 {
			gas = math.MaxUint64 / 2
		}
	} else {
		gas = math.MaxUint64 / 2
	}
	// give an initValue when gasPrice is nil
	var gasPrice *big.Int
	if args.GasPrice == nil {
		gasPrice = new(big.Int).SetUint64(types.DefaultGasPrice)
	} else {
		gasPrice = args.GasPrice.ToBigInt()
		if gasPrice.Sign() == 0 {
			gasPrice = new(big.Int).SetUint64(types.DefaultGasPrice)
		}
	}

	// give an initValue when from is nil
	var from types.Address
	if &args.From == nil || args.From == (types.Address{}) {
		_, addr := wtypes.DefaultTestKey()
		from = types.Address(addr)
	} else {
		from = args.From
	}

	bc, err := repository.NewLatestStateRepository()
	if err != nil {
		return cmn.Bytes{}, fmt.Errorf("new block chain failed")
	}
	// new types.Transaction base on SendTxArgs
	tx := types.NewTransaction(
		bc.GetNonce(*types.TypeConvert(&from)),
		to,
		value,
		gas,
		gasPrice,
		data,
		from,
	)
	result, _, _, err := doCall(tx, blockNr)
	return (cmn.Bytes)(result), err
}

func doCall(tx *craft.Transaction, blockNr types.BlockNumber) ([]byte, uint64, bool, error) {
	bc, err := repository.NewLatestStateRepository()
	if err != nil {
		return nil, 0, true, err
	}
	var block *craft.Block
	if blockNr == types.LatestBlockNumber {
		block = bc.GetCurrentBlock()
	} else {
		height := blockNr.Touint64()
		block, err = bc.GetBlockByHeight(height)
		if err != nil {
			return nil, 0, true, err
		}
	}

	bchash, err := repository.NewRepositoryByBlockHash(block.HeaderHash)
	if err != nil {
		return nil, 0, true, err
	}

	gp := new(common.GasPool).AddGas(uint64(65536))
	result, gas, failed, err, _ := worker.ApplyTransaction(block.Header.Coinbase, block.Header, bchash, tx, gp)
	return result, gas, failed, err
}

//#### eth_gasPrice
//
//Returns the current price per gas in wei.
//
//##### Parameters
//none
//
//##### Returns
//
//`QUANTITY` - integer of the current gas price in wei.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":73}'
//
//// Result
//{
//  "id":73,
//  "jsonrpc": "2.0",
//  "result": "0x09184e72a000" // 10000000000000
//}
//```
//
//***
func GasPrice() (*cmn.Big, error) {
	return (*cmn.Big)(new(big.Int).SetUint64(types.DefaultGasPrice)), nil
}

//#### eth_estimateGas
//
//Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain. Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.
//
//##### Parameters
//
//See [eth_call](#eth_call) parameters, expect that all properties are optional. If no gas limit is specified geth uses the block gas limit from the pending block as an upper bound. As a result the returned estimate might not be enough to executed the call/transaction when the amount of gas is higher than the pending block gas limit.
//
//##### Returns
//
//`QUANTITY` - the amount of gas used.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_estimateGas","params":[{see above}],"id":1}'
//
//// Result
//{
//  "id":1,
//  "jsonrpc": "2.0",
//  "result": "0x5208" // 21000
//}
//```
//
//***
func EstimateGas(args ctypes.SendTxArgs) (cmn.Uint64, error) {
	return cmn.Uint64(21000), nil
}

func Accounts() ([]types.Address, error) {
	addresses := make([]types.Address, 0) // return [] instead of nil if empty
	_, addr := wtypes.DefaultTestKey()
	addresses = append(addresses, types.Address(addr))

	return addresses, nil
}

func Listening() (bool, error) {
	return true, nil // always listening
}

func Version() (string, error) {
	chainId, err := config.GetChainIdFromConfig()
	if err != nil {
		return "1", nil
	}

	id := fmt.Sprint(chainId)
	return id, nil
}
