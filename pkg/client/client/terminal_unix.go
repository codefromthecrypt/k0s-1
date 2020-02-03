// +build linux darwin
// +build !windows

package client

import (
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btwiuse/asciitransport"
	"github.com/containerd/console"
	"k0s.io/k0s/pkg/uuid"
)

func (cl *client) terminalConnect(endpoint string, userinfo *url.Userinfo) {
	log.Println("Press ESC twice to exit. (Some terminals don't supported it)")

	var (
		c    = cl.Config
		conn net.Conn
		err  error
		h    = http.Header{
			"Authorization": {
				"Basic " + base64.StdEncoding.EncodeToString([]byte(userinfo.String())),
			},
		}
	)

	for {
		conn, err = cl.dial(endpoint, h)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	term, err := console.ConsoleFromFile(os.Stdin)
	if err != nil {
		panic(err)
	}
	defer term.Reset()

	if err = term.SetRaw(); err != nil {
		panic(err)
	}

	opts := []asciitransport.Opt{
		asciitransport.WithReader(os.Stdin),
		asciitransport.WithWriter(os.Stdout),
	}

	if c.GetRecord() {
		logname := uuid.New() + ".log"
		logfile, err := os.Create(logname)
		if err != nil {
			panic(err)
		}
		defer func() {
			log.Println("log written to", logname)
		}()
		opts = append(opts, asciitransport.WithLogger(logfile))
	}

	client := asciitransport.Client(conn, opts...)

	// send
	// r
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGWINCH)

		for {
			currentSize, err := term.Size()
			if err != nil {
				log.Println(err)
				continue
			}

			// log.Println(currentSize)
			client.Resize(
				uint(currentSize.Height),
				uint(currentSize.Width),
			)

			switch <-sig {
			case syscall.SIGWINCH:
			}
		}
	}()

	<-client.Done()
}
