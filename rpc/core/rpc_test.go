package core

import (
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/go-kit/kit/log/term"

	"github.com/DSiSc/apigateway/libs/log"
	amino "github.com/tendermint/go-amino"

	server "github.com/DSiSc/apigateway/rpc/lib/server"
)

// Client and Server should work over tcp or unix sockets
const (
	tcpAddr = "tcp://0.0.0.0:47768"

	unixSocket = "/tmp/rpc_test.sock"
	unixAddr   = "unix://" + unixSocket

	websocketEndpoint = "/websocket/endpoint"
)

// Amino codec required to encode/decode everything above.
var RoutesCdc = amino.NewCodec()

var colorFn = func(keyvals ...interface{}) term.FgBgColor {
	for i := 0; i < len(keyvals)-1; i += 2 {
		if keyvals[i] == "socket" {
			if keyvals[i+1] == "tcp" {
				return term.FgBgColor{Fg: term.DarkBlue}
			} else if keyvals[i+1] == "unix" {
				return term.FgBgColor{Fg: term.DarkCyan}
			}
		}
	}
	return term.FgBgColor{}
}

// launch unix and tcp servers
func setup() {
	logger := log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)

	cmd := exec.Command("rm", "-f", unixSocket)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	if err = cmd.Wait(); err != nil {
		panic(err)
	}

	AddTestRoutes()

	tcpLogger := logger.With("socket", "tcp")
	mux := http.NewServeMux()
	server.RegisterRPCFuncs(mux, Routes, RoutesCdc, tcpLogger)
	wm := server.NewWebsocketManager(Routes, RoutesCdc, server.ReadWait(5*time.Second), server.PingPeriod(1*time.Second))
	wm.SetLogger(tcpLogger)
	mux.HandleFunc(websocketEndpoint, wm.WebsocketHandler)
	go func() {
		_, err := server.StartHTTPServer(tcpAddr, mux, tcpLogger, server.Config{})
		if err != nil {
			panic(err)
		}
	}()

	unixLogger := logger.With("socket", "unix")
	mux2 := http.NewServeMux()
	server.RegisterRPCFuncs(mux2, Routes, RoutesCdc, unixLogger)
	wm = server.NewWebsocketManager(Routes, RoutesCdc)
	wm.SetLogger(unixLogger)
	mux2.HandleFunc(websocketEndpoint, wm.WebsocketHandler)
	go func() {
		_, err := server.StartHTTPServer(unixAddr, mux2, unixLogger, server.Config{})
		if err != nil {
			panic(err)
		}
	}()

	// wait for servers to start
	time.Sleep(time.Second * 2)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
