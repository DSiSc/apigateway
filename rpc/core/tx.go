package core

import (
	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
)

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
// | gasPrice  | QUANTITY    | To-Be-Determined  | Option   | Integer of the gasPrice used for each paid gas.                                                      |
// | value     | QUANTITY    | nil               | Option   | Integer of the value sent with this transaction.                                                     |
// | data      | DATA        | nil               | true     | The compiled code of a contract OR the hash of the invoked method signature and encoded parameters.  |
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
func SendTransaction(args ctypes.SendTxArgs) (string, error) {

	return "Not Implement", nil
}
