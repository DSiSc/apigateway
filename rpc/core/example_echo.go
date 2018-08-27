package core

import (
	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
)

type StringArgs struct {
	From string `json:"from"`
}

type ResultEcho struct {
	Value string `json:"value"`
}

func EchoResult(v string) (*ResultEcho, error) {

	return &ResultEcho{v}, nil
}

func EchoResultArgs(v ctypes.StringArgs) (*ResultEcho, error) {

	return &ResultEcho{v.From}, nil
}
