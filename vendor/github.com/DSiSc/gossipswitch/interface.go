package gossipswitch

// GossipSwitchAPI is gossipswitch's public api.
type GossipSwitchAPI interface {
	// NewGossipSwitch create a new switch instance with given filter.
	// filter is used to verify the received message
	NewGossipSwitch(filter SwitchFilter) *GossipSwitch

	// NewGossipSwitchByType create a new switch instance by type.
	// switchType is used to specify the switch type
	NewGossipSwitchByType(switchType SwitchType) (*GossipSwitch, error)

	// InPort get switch's in port by port id, return nil if there is no port with specific id.
	InPort(portId int) *InPort

	// InPort get switch's out port by port id, return nil if there is no port with specific id.
	OutPort(portId int) *OutPort

	// Start start the switch. Once started, switch will receive message from in port, and broadcast to
	// out port
	Start() error

	// Stop stop the switch. Once stopped, switch will stop to receive and broadcast message
	Stop() error

	// IsRunning is used to query switch's current status. Return true if running, otherwise false
	IsRunning() bool
}
