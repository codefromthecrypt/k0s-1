package hub

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"net/url"
	"strconv"

	"github.com/btwiuse/conntroll/pkg/wrap"
	"github.com/btwiuse/gods/maps/linkedhashmap"

	"google.golang.org/grpc"
)

type Agent struct {
	RPCClient      *rpc.Client
	Info           url.Values
	GRPCClientConn chan *grpc.ClientConn
}

type AgentPool struct {
	*linkedhashmap.Map
}

var GlobalAgentPool = &AgentPool{Map: linkedhashmap.New()}

func (p *AgentPool) Del(uuid string) {
	p.Remove(uuid)
}

func (p *AgentPool) Get(uuid string) *Agent {
	v, _ := p.Map.Get(uuid)
	return v.(*Agent)
}

func (p *AgentPool) Add(agent *Agent) {
	p.Put(agent.Info.Get("id"), agent)
}

func (p *AgentPool) Dump() {
	log.Println("[agent pool]")
	for i, v := range p.Values() {
		agent := v.(*Agent)
		uuid := p.Keys()[i].(string)
		fmt.Println(
			fmt.Sprintf("[%s]", strconv.Itoa(i+1)),
			uuid,
			agent.Info,
		)
	}
}

func (p *AgentPool) Has(uuid string) bool {
	_, found := p.Map.Get(uuid)
	return found
}

// we use NewRPCClient over rpc.NewClient(conn)
// so we can remove agent from pool immediately when it is disconnected

/*
                c                           b                  a
          / io.Reader >--->copy()---> io.PipeWriter ===> io.PipeReader = intercepted io.Reader \
net.Conn  - io.Writer \                                                                        wrap.ReadWriteCloser() - rpc.NewClient - *rpc.Client
          \ io.Closer - io.WriteCloser ---------------------------------------------------------
*/
func (agent *Agent) MakeInterceptedRPCClient(c io.ReadWriteCloser) {
	a, b := io.Pipe()
	go func() {
		defer agent.onClose()
		if _, err := io.Copy(b, c); err != nil {
			log.Println(err)
		}
	}()
	agent.RPCClient = rpc.NewClient(wrap.WrapReadWriteCloser(a, c))
}

// onclose is called when agent goes offline
func (agent *Agent) onClose() {
	// TODO: remove Dump
	// panic: runtime error: index out of range [3] with length 3
	// defer GlobalAgentPool.Dump()
	log.Println("disconnected:", agent.Info.Get("id"))
	GlobalAgentPool.Del(agent.Info.Get("id"))
}
