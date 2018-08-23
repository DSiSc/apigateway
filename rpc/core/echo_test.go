package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
	client "github.com/DSiSc/apigateway/rpc/lib/client"
)

func TestEchoViaHTTP(t *testing.T) {

	cl := client.NewJSONRPCClient(tcpAddr)

	testVal := "abcd"

	params := map[string]interface{}{
		"arg": testVal,
	}

	result := new(ResultEcho)
	_, err := cl.Call("echo", params, result)

	require.Nil(t, err)
	assert.Equal(t, result.Value, testVal)

}

func TestEchoArgsViaHTTP(t *testing.T) {

	cl := client.NewJSONRPCClient(tcpAddr)

	testVal := "abcd"

	param_args := ctypes.StringArgs{
		From: testVal,
	}

	params := map[string]interface{}{
		"arg": param_args,
	}

	result := new(ResultEcho)
	_, err := cl.Call("echo_args", params, result)

	require.Nil(t, err)
	assert.Equal(t, result.Value, testVal)
}
