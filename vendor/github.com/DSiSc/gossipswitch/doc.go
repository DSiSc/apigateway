/*
Package gossipswitch implements the gossip switch(receive the message from InPort, filter the received message,
and then broadcast the message to OutPort).

The gossipswitch package implements one complete gossip switch, include: switch InPort(receive message), switch
OutPort(send message),two common message filter(TxFilter: filter transaction message, BlockFilter:filter block
message).
*/
package gossipswitch
