package core

import (
	"encoding/json"
	"errors"
	coretypes "github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/apigateway/rpc/lib/server"
	"github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/DSiSc/blockchain"
	bcconf "github.com/DSiSc/blockchain/config"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/monkey"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	"math/big"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func mockBlock() *types.Block {
	header := &types.Header{}
	block := &types.Block{
		HeaderHash: types.Hash{
			0xbd, 0x79, 0x1d, 0x4a, 0xf9, 0x64, 0x8f, 0xc3, 0x7f, 0x94, 0xeb, 0x36, 0x53, 0x19, 0xf6, 0xd0,
			0xa9, 0x78, 0x9f, 0x9c, 0x22, 0x47, 0x2c, 0xa7, 0xa6, 0x12, 0xa9, 0xca, 0x4, 0x13, 0xc1, 0x4,
		},
		Header: header,
	}
	return block
}

// mock receipts
func mockReceipts() []*types.Receipt {
	txHash := types.Hash{
		0xbd, 0x79, 0x1d, 0x4a, 0xf9, 0x64, 0x8f, 0xc3, 0x7f, 0x94, 0xeb, 0x36, 0x53, 0x19, 0xf6, 0xd0,
		0xa9, 0x78, 0x9f, 0x9c, 0x22, 0x47, 0x2c, 0xa7, 0xa6, 0x12, 0xa9, 0xca, 0x4, 0x13, 0xc1, 0x4,
	}
	log := &types.Log{
		TxHash: txHash,
	}
	receipt := types.Receipt{
		Status: 1,
		Logs:   []*types.Log{log},
	}
	return []*types.Receipt{&receipt}
}

func mockTx() *types.Transaction {
	return &types.Transaction{
		Data: types.TxData{
			Price: big.NewInt(1),
		},
	}
}

// test subscribe
func TestSubscribe(t *testing.T) {
	assert := assert.New(t)
	defer monkey.UnpatchAll()
	monkey.Patch(NewHeaders, func(rpctypes.WSRPCContext) (string, error) { return "newHeads", nil })
	monkey.Patch(Logs, func(rpctypes.WSRPCContext, json.RawMessage) (string, error) { return "logs", nil })
	monkey.Patch(NewPendingTransactions, func(rpctypes.WSRPCContext) (string, error) { return "newPendingTransactions", nil })

	result, err := Subscribe(rpctypes.WSRPCContext{}, strArrayToJsonRawArray("\"newHeads\""))
	assert.Nil(err)
	assert.Equal("newHeads", result)

	result, err = Subscribe(rpctypes.WSRPCContext{}, strArrayToJsonRawArray("\"logs\"", "\"test\""))
	assert.Nil(err)
	assert.Equal("logs", result)

	result, err = Subscribe(rpctypes.WSRPCContext{}, strArrayToJsonRawArray("\"newPendingTransactions\""))
	assert.Nil(err)
	assert.Equal("newPendingTransactions", result)
}

// test subscribe new headers event
func TestSubscribeNewHeaders(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conn := new(websocket.Conn)
	monkey.PatchInstanceMethod(reflect.TypeOf(conn), "RemoteAddr", func(*websocket.Conn) net.Addr {
		return &net.IPNet{IP: net.ParseIP("127.0.0.1")}
	})

	ec := newTestEvent()
	wsc := rpcserver.NewWSConnection(conn, Routes, amino.NewCodec(), rpcserver.EventSubscriber(ec))
	twsc := newTestConnection(wsc)
	wctx := rpctypes.WSRPCContext{WSRPCConnection: twsc}

	result, err := Subscribe(wctx, strArrayToJsonRawArray("\"newHeads\""))
	assert.Nil(err)
	assert.NotNil(result)

	block := mockBlock()
	ec.Notify(types.EventBlockCommitted, block)
	timer := time.NewTicker(time.Second)
	select {
	case returnResp := <-twsc.WriteChan:
		notify := make(map[string]interface{})
		notify["subscription"] = result
		notify["result"] = block.Header
		b, _ := json.Marshal(notify)
		rb, _ := returnResp.Params.MarshalJSON()
		assert.Equal(b, rb)
	case <-timer.C:
		assert.Nil(errors.New("Failed to subscribe newHeads event"))
	}

}

// test subscribe new transaction event
func TestNewPendingTransactions(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conn := new(websocket.Conn)
	monkey.PatchInstanceMethod(reflect.TypeOf(conn), "RemoteAddr", func(*websocket.Conn) net.Addr {
		return &net.IPNet{IP: net.ParseIP("127.0.0.1")}
	})

	ec := newTestEvent()
	wsc := rpcserver.NewWSConnection(conn, Routes, amino.NewCodec(), rpcserver.EventSubscriber(ec))
	twsc := newTestConnection(wsc)
	wctx := rpctypes.WSRPCContext{WSRPCConnection: twsc}

	result, err := Subscribe(wctx, strArrayToJsonRawArray("\"newPendingTransactions\""))
	assert.Nil(err)
	assert.NotNil(result)

	tx := mockTx()
	ec.Notify(types.EventAddTxToTxPool, coretypes.TxHash(tx))
	timer := time.NewTicker(time.Second)
	select {
	case returnResp := <-twsc.WriteChan:
		notify := make(map[string]interface{})
		notify["subscription"] = result
		notify["result"] = tx.Hash.Load()
		b, _ := json.Marshal(notify)
		rb, _ := returnResp.Params.MarshalJSON()
		assert.Equal(b, rb)
	case <-timer.C:
		assert.Nil(errors.New("Failed to subscribe newPendingTransactions event"))
	}
}

// test subscribe contract execution log
func TestLogs(t *testing.T) {
	defer monkey.UnpatchAll()
	assert := assert.New(t)
	conn := new(websocket.Conn)
	monkey.PatchInstanceMethod(reflect.TypeOf(conn), "RemoteAddr", func(*websocket.Conn) net.Addr {
		return &net.IPNet{IP: net.ParseIP("127.0.0.1")}
	})

	ec := newTestEvent()
	wsc := rpcserver.NewWSConnection(conn, Routes, amino.NewCodec(), rpcserver.EventSubscriber(ec))
	twsc := newTestConnection(wsc)
	wctx := rpctypes.WSRPCContext{WSRPCConnection: twsc}

	result, err := Subscribe(wctx, strArrayToJsonRawArray("\"logs\"", "{}"))
	assert.Nil(err)
	assert.NotNil(result)

	block := mockBlock()
	blockchain.InitBlockChain(bcconf.BlockChainConfig{PluginName: blockchain.PLUGIN_MEMDB}, ec)
	bChain, _ := blockchain.NewLatestStateBlockChain()
	receipts := mockReceipts()
	bChain.WriteBlockWithReceipts(block, receipts)

	timer := time.NewTicker(1000 * time.Second)
	select {
	case returnResp := <-twsc.WriteChan:
		notify := make(map[string]interface{})
		notify["subscription"] = result
		notify["result"] = receipts[0].Logs
		b, _ := json.Marshal(notify)
		rb, _ := returnResp.Params.MarshalJSON()
		assert.Equal(b, rb)
	case <-timer.C:
		assert.Nil(errors.New("Failed to subscribe newPendingTransactions event"))
	}
}

func strArrayToJsonRawArray(strs ...string) []json.RawMessage {
	raw := make([]json.RawMessage, 0)
	for _, str := range strs {
		raw = append(raw, []byte(str))
	}
	return raw
}

type testEvent struct {
	m           sync.RWMutex
	Subscribers map[types.EventType]map[types.Subscriber]types.EventFunc
}

func newTestEvent() types.EventCenter {
	return &testEvent{
		Subscribers: make(map[types.EventType]map[types.Subscriber]types.EventFunc),
	}
}

//  adds a new subscriber to testEvent.
func (e *testEvent) Subscribe(eventType types.EventType, eventFunc types.EventFunc) types.Subscriber {
	e.m.Lock()
	defer e.m.Unlock()

	sub := make(chan interface{})
	_, ok := e.Subscribers[eventType]
	if !ok {
		e.Subscribers[eventType] = make(map[types.Subscriber]types.EventFunc)
	}
	e.Subscribers[eventType][sub] = eventFunc

	return sub
}

// UnSubscribe removes the specified subscriber
func (e *testEvent) UnSubscribe(eventType types.EventType, subscriber types.Subscriber) (err error) {
	e.m.Lock()
	defer e.m.Unlock()

	subEvent, ok := e.Subscribers[eventType]
	if !ok {
		err = errors.New("event type not exist")
		return
	}

	delete(subEvent, subscriber)
	close(subscriber)

	return
}

// Notify subscribers that Subscribe specified event
func (e *testEvent) Notify(eventType types.EventType, value interface{}) (err error) {

	e.m.RLock()
	defer e.m.RUnlock()

	subs, ok := e.Subscribers[eventType]
	if !ok {
		err = errors.New("event type not register")
		return
	}

	switch value.(type) {
	case error:
		log.Error("Receive errors is [%v].", value)
	}
	log.Info("Receive eventType is [%d].", eventType)

	for _, event := range subs {
		go e.NotifySubscriber(event, value)
	}
	return nil
}

func (e *testEvent) NotifySubscriber(eventFunc types.EventFunc, value interface{}) {
	if eventFunc == nil {
		return
	}

	// invoke subscriber event func
	eventFunc(value)

}

//Notify all event subscribers
func (e *testEvent) NotifyAll() (errs []error) {
	e.m.RLock()
	defer e.m.RUnlock()

	for eventType, _ := range e.Subscribers {
		if err := e.Notify(eventType, nil); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// unsubscribe all event and subscriber elegant
func (e *testEvent) UnSubscribeAll() {
	e.m.Lock()
	defer e.m.Unlock()
	for eventtype, _ := range e.Subscribers {
		subs, ok := e.Subscribers[eventtype]
		if !ok {
			continue
		}
		for subscriber, _ := range subs {
			delete(subs, subscriber)
			close(subscriber)
		}
	}
	// TODO: open it when txswitch and blkswith stop complete
	//e.Subscribers = make(map[types.EventType]map[types.Subscriber]types.EventFunc)
	return
}

type TestConnection struct {
	WriteChan chan rpctypes.RPCResponse
	wc        rpctypes.WSRPCConnection
}

func newTestConnection(wc rpctypes.WSRPCConnection) *TestConnection {
	return &TestConnection{
		wc:        wc,
		WriteChan: make(chan rpctypes.RPCResponse),
	}
}

// GetRemoteAddr returns the remote address of the underlying connection.
// It implements WSRPCConnection
func (wsc *TestConnection) GetRemoteAddr() string {
	return wsc.wc.GetRemoteAddr()
}

// GetEventSubscriber implements WSRPCConnection by returning event subscriber.
func (wsc *TestConnection) GetEventSubscriber() rpctypes.EventSubscriber {
	return wsc.wc.GetEventSubscriber()
}

// WriteRPCResponse pushes a response to the writeChan, and blocks until it is accepted.
// It implements WSRPCConnection. It is Goroutine-safe.
func (wsc *TestConnection) WriteRPCResponse(resp rpctypes.RPCResponse) {
	wsc.WriteChan <- resp
}

// TryWriteRPCResponse attempts to push a response to the writeChan, but does not block.
// It implements WSRPCConnection. It is Goroutine-safe
func (wsc *TestConnection) TryWriteRPCResponse(resp rpctypes.RPCResponse) bool {
	return wsc.wc.TryWriteRPCResponse(resp)
}

// Codec returns an amino codec used to decode parameters and encode results.
// It implements WSRPCConnection.
func (wsc *TestConnection) Codec() *amino.Codec {
	return wsc.wc.Codec()
}
