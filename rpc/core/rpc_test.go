package core

import (
	"bytes"
	"net/http"

	"github.com/DSiSc/apigateway/log"
	rs "github.com/DSiSc/apigateway/rpc/lib/server"
	amino "github.com/tendermint/go-amino"
)

func testMux() *http.ServeMux {

	cdc := amino.NewCodec()
	mux := http.NewServeMux()
	buf := new(bytes.Buffer)
	logger := log.NewTMLogger(buf)
	rs.RegisterRPCFuncs(mux, Routes, cdc, logger)

	return mux
}

func statusOK(code int) bool { return code >= 200 && code <= 299 }
