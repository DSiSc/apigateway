package core

import (
	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
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
