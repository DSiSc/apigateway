package apigateway

import (
	"net"
	"net/http"
	"os"
	"time"

	cmn "github.com/DSiSc/apigateway/common"
	"github.com/DSiSc/apigateway/log"
	//	ctypes "github.com/DSiSc/apigateway/rpc/core/types"
	"github.com/go-kit/kit/log/term"
	amino "github.com/tendermint/go-amino"

	rpccore "github.com/DSiSc/apigateway/rpc/core"
	rpcserver "github.com/DSiSc/apigateway/rpc/lib/server"
	craftlog "github.com/DSiSc/craft/log"
)

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

var logger = log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)

func StartRPC(listenAddr string) ([]net.Listener, error) {

	listenAddrs := cmn.SplitAndTrim(listenAddr, ",", " ")
	coreCodec := amino.NewCodec()
	// TODO(peerlink): let's see wire.go
	//ctypes.RegisterAmino(coreCodec)

	// we may expose the rpc over both a unix and tcp socket
	listeners := make([]net.Listener, len(listenAddrs))
	for i, listenAddr := range listenAddrs {
		mux := http.NewServeMux()
		rpcLogger := logger.With("module", "rpc-server")
		wm := rpcserver.NewWebsocketManager(rpccore.Routes, coreCodec, rpcserver.ReadWait(5*time.Second))
		// TODO(peerlink): rpcserver get eventBus from input vars.
		//rpcserver.EventSubscriber(n.eventBus))
		wm.SetLogger(rpcLogger.With("protocol", "websocket"))
		mux.HandleFunc("/websocket", wm.WebsocketHandler)
		rpcserver.RegisterRPCFuncs(mux, rpccore.Routes, coreCodec, rpcLogger)
		listener, err := rpcserver.StartHTTPServer(
			listenAddr,
			mux,
			rpcLogger,
			// TODO(peerlink): rpcserver.Config get MaxOpenConnections from input vars.
			//rpcserver.Config{MaxOpenConnections: maxOpenConnections},
			rpcserver.Config{},
		)
		if err != nil {
			return nil, err
		}
		listeners[i] = listener
	}
	return listeners, nil
}

// StopRPC stop RPC server
func StopRPC(rpcListeners []net.Listener) error {

	for _, l := range rpcListeners {
		//logger.Info("Closing rpc listener", "listener", l)
		craftlog.InfoKV("Closing rpc listener", map[string]interface{}{"listener": l})
		if err := l.Close(); err != nil {
			//logger.Error("Error closing listener", "listener", l, "err", err)
			craftlog.ErrorKV("Error closing listener", map[string]interface{}{"listener": l, "err": err})
			return err
		}
	}
	return nil
}
