package core

import (
	"encoding/json"
	"fmt"
	"github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/DSiSc/craft/log"
	crafttypes "github.com/DSiSc/craft/types"
	"github.com/DSiSc/repository"
	"github.com/pkg/errors"
	"math/big"
)

const (
	NewHeadersEvent             = "newHeads"
	LogsEvent                   = "logs"
	NewPendingTransactionsEvent = "newPendingTransactions"
	SyncingEvent                = "syncing"
)

// FilterCriteria contains options for contract log filtering.
type FilterCriteria struct {
	BlockHash *crafttypes.Hash     `json:"block_hash"` // used by eth_getLogs, return logs only from block with this hash
	FromBlock *big.Int             `json:"from_block"` // beginning of the queried range,	 nil means genesis block
	ToBlock   *big.Int             `json:"to_block"`   // end of the range, nil means latest block
	Addresses []crafttypes.Address `json:"addresses"`  // restricts matches to events created by specific contracts

	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Examples:
	// {} or nil          matches any topic list
	// {{A}}              matches topic A in first position
	// {{}, {B}}          matches any topic in first position AND B in second position
	// {{A}, {B}}         matches topic A in first position AND B in second position
	// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
	Topics [][]crafttypes.Hash `json:"topics"`
}

//#### eth_subscribe
//
//Subscribe for events(newHeads/logs/newPendingTransactions) via WebSocket.
//
//
//##### Parameters
//
//1. `TAG` - subscription name `"newHeads"`, `"logs"` or `"newPendingTransactions"`（newHeads: new header is appended to the chain; logs: new logs are included in new blocks; newPendingTransactions: new transactions are added to the pending state and are signed with a key that is available in the node）.
//2. `Object` - the transaction index position.
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
//subscription id
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"eth_subscribe","params":["newHeads"],"id":1}'
//```
//
//// Result
//{
//"id":1,
//"jsonrpc": "2.0",
//"result": "0x919d38fa5c395fa0f677e6554eef74fc7"
//}
//***
//Subscribe for events via WebSocket.
func Subscribe(wsCtx rpctypes.WSRPCContext, rawMsg []json.RawMessage) (string, error) {
	var subscribeMethod string
	if err := json.Unmarshal(rawMsg[0], &subscribeMethod); err != nil {
		log.Error("Unable to parse subscription method: %v\n", err)
		return "", errors.New("Unknown subscription method")
	}
	switch subscribeMethod {
	case NewHeadersEvent:
		return NewHeaders(wsCtx)
	case LogsEvent:
		return Logs(wsCtx, rawMsg[1])
	case NewPendingTransactionsEvent:
		return NewPendingTransactions(wsCtx)
	case SyncingEvent:
		//TODO
	}
	return "", errors.New("Unknown subscription method")
}

//NewHeaders subscribe new headers event
func NewHeaders(wsCtx rpctypes.WSRPCContext) (string, error) {
	subscription, err := wsCtx.GetEventSubscriber().Subscribe(crafttypes.EventBlockWritten, crafttypes.EventBlockCommitted)
	if err != nil {
		return "", err
	}
	// send notification when received a new event
	go func(wsCtx rpctypes.WSRPCContext, sub *rpctypes.Subscription) {
		for {
			select {
			case event := <-sub.EventChan():
				fmt.Println("sssssssssss")
				if block, ok := event.(*crafttypes.Block); ok {
					if resp, err := rpctypes.NewJsonEventNotifyResponse(sub.ID, block.Header); err == nil {
						wsCtx.WriteRPCResponse(resp)
					}
				}
			case <-sub.QuitChan():
				return
			}
		}
	}(wsCtx, subscription)
	return subscription.ID, nil
}

//Logs subscribe new logs event
func Logs(wsCtx rpctypes.WSRPCContext, rawMsg json.RawMessage) (string, error) {
	crit := new(FilterCriteria)
	err := json.Unmarshal(rawMsg, crit)
	if err != nil {
		return "", errors.New("Failed to parse filter criteria ")
	}
	subscription, err := wsCtx.GetEventSubscriber().Subscribe(crafttypes.EventBlockWritten, crafttypes.EventBlockCommitted)
	if err != nil {
		return "", err
	}
	// send notification when received a new event
	go func(wsCtx rpctypes.WSRPCContext, sub *rpctypes.Subscription, crit *FilterCriteria) {
		for {
			select {
			case event := <-sub.EventChan():
				if block, ok := event.(*crafttypes.Block); ok {
					headerHash := types.HeaderHash(block)
					blockChain, err := repository.NewLatestStateRepository()
					if err != nil {
						log.Warn("Failed to get latest blockchain, as: %v ", err)
						continue
					}
					receipts := blockChain.GetReceiptByBlockHash(headerHash)
					if nil != receipts {
						logsList := make([][]*crafttypes.Log, len(receipts))
						for i, receipt := range receipts {
							logsList[i] = receipt.Logs
						}
						var unfiltered []*crafttypes.Log
						for _, logs := range logsList {
							unfiltered = append(unfiltered, logs...)
						}
						logs := filterLogs(unfiltered, crit.FromBlock, crit.ToBlock, crit.Addresses, crit.Topics)
						if resp, err := rpctypes.NewJsonEventNotifyResponse(sub.ID, logs); err == nil {
							wsCtx.WriteRPCResponse(resp)
						}
					}
				}
			case <-sub.QuitChan():
				return
			}
		}
	}(wsCtx, subscription, crit)
	return subscription.ID, nil
}

//NewPendingTransactions subscribe new transactions event
func NewPendingTransactions(wsCtx rpctypes.WSRPCContext) (string, error) {
	subscription, err := wsCtx.GetEventSubscriber().Subscribe(crafttypes.EventAddTxToTxPool)
	if err != nil {
		return "", err
	}
	// send notification when received a new event
	go func(wsCtx rpctypes.WSRPCContext, sub *rpctypes.Subscription) {
		for {
			select {
			case event := <-sub.EventChan():
				if tx, ok := event.(*crafttypes.Transaction); ok {
					if resp, err := rpctypes.NewJsonEventNotifyResponse(sub.ID, tx.Hash.Load()); err == nil {
						wsCtx.WriteRPCResponse(resp)
					}
				}
			case <-sub.QuitChan():
				return
			}
		}
	}(wsCtx, subscription)
	return subscription.ID, nil
}

//Subscribe for events via WebSocket.
func UnSubscribe(wsCtx rpctypes.WSRPCContext, subID string) (bool, error) {
	err := wsCtx.GetEventSubscriber().Unsubscribe(subID)
	return nil == err, err
}

// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*crafttypes.Log, fromBlock, toBlock *big.Int, addresses []crafttypes.Address, topics [][]crafttypes.Hash) []*crafttypes.Log {
	var ret []*crafttypes.Log
Logs:
	for _, log := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > log.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < log.BlockNumber {
			continue
		}

		if len(addresses) > 0 && !includes(addresses, log.Address) {
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(log.Topics) {
			continue Logs
		}
		for i, sub := range topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if log.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, log)
	}
	return ret
}

func includes(addresses []crafttypes.Address, a crafttypes.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}
	return false
}
