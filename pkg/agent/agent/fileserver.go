package agent

import (
	"net"
	"net/http"

	types "k0s.io/k0s/pkg/agent"
)

func StartFileServer(c types.Config) chan net.Conn {
	var (
		fsListener = NewLys()
		handler    = LoggingMiddleware(http.FileServer(http.Dir("/")))
		fileServer = &http.Server{Handler: handler}
	)
	go fileServer.Serve(fsListener)
	return fsListener.Conns
}
