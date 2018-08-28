package gossipswitch

import (
	"errors"
	"sync"
	"sync/atomic"
)

// SwitchType switch type
type SwitchType int

const (
	TxSwitch SwitchType = iota
	BlockSwitch
)

// SwitchMsg is the message that can be dealt with by GossipSwitch.
type SwitchMsg interface {
}

// Filter is used to verify SwitchMsg
type SwitchFilter interface {
	Verify(msg SwitchMsg) error
}

// common const value
const (
	LocalInPortId   = 0 //Local InPort ID, receive the message from local
	RemoteInPortId  = 1 //Remote InPort ID, receive the message from remote
	LocalOutPortId  = 0 //Local OutPort ID
	RemoteOutPortId = 1 //Remote OutPort ID
)

// GossipSwitch is the implementation of gossip switch.
// for gossipswitch, if a validated message is received, it will be broadcasted,
// otherwise it will be dropped.
type GossipSwitch struct {
	switchMtx sync.Mutex
	filter    SwitchFilter
	inPorts   map[int]*InPort
	outPorts  map[int]*OutPort
	isRunning uint32 // atomic
}

// NewGossipSwitch create a new switch instance with given filter.
// filter is used to verify the received message
func NewGossipSwitch(filter SwitchFilter) *GossipSwitch {
	sw := &GossipSwitch{
		filter:   filter,
		inPorts:  make(map[int]*InPort),
		outPorts: make(map[int]*OutPort),
	}
	sw.initPort()
	return sw
}

// NewGossipSwitchByType create a new switch instance by type.
// switchType is used to specify the switch type
func NewGossipSwitchByType(switchType SwitchType) (*GossipSwitch, error) {
	var filter SwitchFilter
	switch switchType {
	case TxSwitch:
		filter = NewTxFilter()
	case BlockSwitch:
		filter = NewBlockFilter()
	default:
		return nil, errors.New("unsupported switch type")
	}
	sw := &GossipSwitch{
		filter:   filter,
		inPorts:  make(map[int]*InPort),
		outPorts: make(map[int]*OutPort),
	}
	sw.initPort()
	return sw, nil
}

// init switch's InPort and OutPort
func (sw *GossipSwitch) initPort() {
	sw.inPorts[LocalInPortId] = newInPort()
	sw.inPorts[RemoteInPortId] = newInPort()
	sw.outPorts[LocalOutPortId] = newOutPort()
	sw.outPorts[RemoteOutPortId] = newOutPort()
}

// InPort get switch's in port by port id, return nil if there is no port with specific id.
func (sw *GossipSwitch) InPort(portId int) *InPort {
	return sw.inPorts[portId]
}

// InPort get switch's out port by port id, return nil if there is no port with specific id.
func (sw *GossipSwitch) OutPort(portId int) *OutPort {
	return sw.outPorts[portId]
}

// Start start the switch. Once started, switch will receive message from in port, and broadcast to
// out port
func (sw *GossipSwitch) Start() error {
	if atomic.CompareAndSwapUint32(&sw.isRunning, 0, 1) {
		for _, inPort := range sw.inPorts {
			go sw.receiveRoutine(inPort)
		}
		return nil
	}
	return errors.New("switch already started")
}

// Stop stop the switch. Once stopped, switch will stop to receive and broadcast message
func (sw *GossipSwitch) Stop() error {
	if atomic.CompareAndSwapUint32(&sw.isRunning, 1, 0) {
		return nil
	}
	return errors.New("switch already stopped")
}

// IsRunning is used to query switch's current status. Return true if running, otherwise false
func (sw *GossipSwitch) IsRunning() bool {
	return atomic.LoadUint32(&sw.isRunning) == 1
}

// listen to receive message from the in port
func (sw *GossipSwitch) receiveRoutine(inPort *InPort) {
	for {
		select {
		case msg := <-inPort.read():
			sw.onRecvMsg(msg)
		}

		//check switch status
		if !sw.IsRunning() {
			break
		}
	}
}

// deal with the received message.
func (sw *GossipSwitch) onRecvMsg(msg SwitchMsg) {
	if err := sw.filter.Verify(msg); err == nil {
		sw.broadCastMsg(msg)
	}
}

// broadcast the validated message to all out ports.
func (sw *GossipSwitch) broadCastMsg(msg SwitchMsg) error {
	for _, outPort := range sw.outPorts {
		go outPort.write(msg)
	}
	return nil
}
