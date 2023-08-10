package pubsub

import (
	"fmt"
	"net"

	"github.com/panupakm/miniredis/payload"
)

type PubSub struct {
	connmap map[string][]net.Conn
}

type PubSuber interface {
	Sub(topic string, conn net.Conn)
	Pub(topic string, typ payload.ValueType, buff []byte, conn net.Conn)
	Unsub(conn net.Conn)
	IsSub(topic string, conn net.Conn) bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		connmap: make(map[string][]net.Conn),
	}
}

func (ps *PubSub) Sub(topic string, conn net.Conn) {
	ps.connmap[topic] = append(ps.connmap[topic], conn)
}

func (ps *PubSub) IsSub(topic string) bool {
	conns, ok := ps.connmap[topic]
	return ok && len(conns) > 0
}

func (ps *PubSub) Pub(topic string, typ payload.ValueType, buff []byte, conn net.Conn) {
	conns := ps.connmap[topic]
	for _, con := range conns {
		if con == conn {
			continue
		}

		msg := payload.NewSubMsg(typ, buff)
		_, err := msg.WriteTo(con)
		if err != nil {
			fmt.Println("pub error:", err)
		}
	}
}

func (ps *PubSub) Unsub(conn net.Conn) {
	for topic, conns := range ps.connmap {
		for i, c := range conns {
			if c == conn {
				ps.connmap[topic] = append(ps.connmap[topic][:i], ps.connmap[topic][i+1:]...)
				break
			}
		}
	}
}
