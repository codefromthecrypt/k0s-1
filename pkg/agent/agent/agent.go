package agent

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os/exec"
	"path"

	types "github.com/btwiuse/conntroll/pkg/agent"
	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

var (
	_ types.Agent = (*agent)(nil)
)

type agent struct {
	*errgroup.Group
	types.Config
	// types.RPC

	types.GRPCServer
	// grpcln chan<- net.Conn
	id   string
	name string
}

func NewAgent(c types.Config) types.Agent {
	eg, _ := errgroup.WithContext(context.Background())
	id := c.GetID()
	name := c.GetName()
	if c.GetVerbose() {
		log.Println("new agent", id, name)
	}
	shell := "bash"
	if _, err := exec.LookPath(shell); err != nil {
		shell = "sh"
	}
	grpcServer := StartGRPCServer(c.GetCmd())

	return &agent{
		Group:      eg,
		Config:     c,
		GRPCServer: grpcServer,
		id:         id,
		name:       name,
	}
}

func (ag *agent) Accept() (net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)

	conn, err = ag.Dial("/api/grpc?id=" + ag.GetID())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (ag *agent) Dial(p string) (conn net.Conn, err error) {
	var (
		c = ag.Config
		u string
	)

	switch c.GetScheme() {
	case "http":
		u = "ws://" + path.Join(c.GetAddr(), p)
	case "https":
		u = "wss://" + path.Join(c.GetAddr(), p)
	}

	wsconn, _, err := websocket.Dial(context.Background(), u, nil)

	if err != nil {
		return nil, err
	}

	return websocket.NetConn(context.Background(), wsconn, websocket.MessageBinary), nil
}

func (ag *agent) AgentRegister(conn net.Conn) (types.RPC, error) {
	_, err := io.WriteString(conn, ag.Config.String())
	if err != nil {
		return nil, err
	}

	return NewRPC(conn), nil
}

func (ag *agent) Serve(rpc types.RPC) error {
	for {
		select {
		case f := <-rpc.Actions():
			go f(ag)
		case <-rpc.Done():
			break
		}
	}
	return errors.New("yrpc connection closed")
}

func (ag *agent) ConnectAndServe() error {
	conn, err := ag.Dial("/api/rpc")
	if err != nil {
		return err
	}

	rpc, err := ag.AgentRegister(conn)
	if err != nil {
		return err
	}

	log.Println("Connect ok")

	err = ag.Serve(rpc)
	if err != nil {
		log.Println("Serve:", err)
		return err
	}

	return nil
}
