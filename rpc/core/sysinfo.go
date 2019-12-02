package core

import (
	"fmt"
	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/repository"
)

func ChannelInfo() ([]ctypes.ChannelInfo, error) {
	channels := make([]ctypes.ChannelInfo, 0) // return [] instead of nil if empty
	channelInfo := ctypes.ChannelInfo{
		Name:      "justitia-chan1",
		ChannelId: "justitia-chan1",
	}
	channels = append(channels, channelInfo)
	return channels, nil
}

func NodeInfo() ([]ctypes.NodeInfo, error) {
	nodeInfos := make([]ctypes.NodeInfo, 0) // return [] instead of nil if empty
	NodeInfo := ctypes.NodeInfo{
		HostName: "justitia-node1",
		Url:      "127.0.0.1:47768",
		Genesis:  "justitia-chan1",
	}
	nodeInfos = append(nodeInfos, NodeInfo)
	return nodeInfos, nil
}

//#### net_sysContract
//
//Get system contracts info.
//
//
//##### Parameters
//
//
//##### Returns
//
//`Map` - system contract info:
//
//- `Key`: `DATA` - contract name.
//- `Value`: `DATA` - contract address.
//
//##### Example
//```js
//// Request
//curl -X POST --data '{"jsonrpc":"2.0","method":"net_sysContract","id":1}'
//
//// Result
//{
//  "jsonrpc": "2.0",
//  "id": "",
//  "result": {
//    "DposPbft": "bd770416a3345f91e4b34576cb804a576fa48eb1",
//    "MspContract": "02370416a3556f91e4a6d576cb804a576fa483da"
//  }
//}
//```
//***
func SystemContract() (map[string]string, error) {
	bc, _ := repository.NewLatestStateRepository()
	sysContracts := make(map[string]string)
	if dposContract, err := bc.Get([]byte(types.DposBftVotingContract)); err == nil {
		sysContracts[types.DposBftVotingContract] = fmt.Sprintf("%x", dposContract)
	}
	if mspContract, err := bc.Get([]byte(types.MspContract)); err == nil {
		sysContracts[types.MspContract] = fmt.Sprintf("%x", mspContract)
	}
	return sysContracts, nil
}
