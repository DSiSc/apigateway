package core

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DSiSc/apigateway/log"
	rs "github.com/DSiSc/apigateway/rpc/lib/server"
	ltypes "github.com/DSiSc/apigateway/rpc/lib/types"
	"github.com/tendermint/go-amino"
)

type Requestdata struct {
	payload    string
	wantErr    string
	wantRetrun string
}

func testMux() *http.ServeMux {

	cdc := amino.NewCodec()
	mux := http.NewServeMux()
	buf := new(bytes.Buffer)
	logger := log.NewTMLogger(buf)
	rs.RegisterRPCFuncs(mux, Routes, cdc, logger)

	return mux
}

func statusOK(code int) bool { return code >= 200 && code <= 299 }

func doRpcTest(t *testing.T, tests []*Requestdata) {
	mux := testMux()
	for i, tt := range tests {
		req, _ := http.NewRequest("POST", "http://localhost/", strings.NewReader(tt.payload))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		// --------------
		// Test Response
		res := rec.Result()
		// Always expecting back a JSONRPCResponse
		assert.True(t, statusOK(res.StatusCode), "#%d: should always return 2XX", i)
		blob, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("#%d: err reading body: %v", i, err)
			continue
		}
		// ----------------
		// Test reponse
		recv := new(ltypes.RPCResponse)
		json.Unmarshal(blob, recv)

		b, _ := json.Marshal(recv)
		assert.Equal(t, tt.wantRetrun, string(b))
	}
}
