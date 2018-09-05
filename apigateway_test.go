package apigateway

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

var (
	tcpAddr = "tcp://0.0.0.0:47768"
)

func TestStartAndStopRPC(t *testing.T) {

	cases := []struct {
		input   string
		want    string
		wanterr string
	}{
		{
			input:   tcpAddr,
			want:    "[::]:47768",
			wanterr: "Error",
		},
	}

	for _, tcase := range cases {
		actualListenr, err := StartRPC(tcase.input)
		require.Nil(t, err)
		assert.Equal(t, tcase.want, actualListenr[0].Addr().String(), "Listener should be same")

		// Test stopRPC
		errStop := StopRPC(actualListenr)
		require.Nil(t, errStop)
	}
}
