package pubsub

import (
	"fmt"
	"net"

	"github.com/panupakm/miniredis/lib/payload"
)

type PubSub struct {
	Map map[string][]net.Conn
}

func NewPubSub() *PubSub {
	return &PubSub{
		Map: make(map[string][]net.Conn),
	}
}

func (ps *PubSub) Sub(topic string, conn net.Conn) {
	ps.Map[topic] = append(ps.Map[topic], conn)
}

func (ps *PubSub) Pub(topic string, typ payload.ValueType, buff []byte, conn net.Conn) {
	fmt.Printf("Pub: %s\n", string(buff))
	conns := ps.Map[topic]
	for _, con := range conns {
		if con == conn {
			continue
		}

		msg := payload.NewSubMsg(typ, buff)
		_, err := msg.WriteTo(con)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (ps *PubSub) UnsubConnection(conn net.Conn) {
	for topic, conns := range ps.Map {
		for i, c := range conns {
			if c == conn {
				ps.Map[topic] = append(ps.Map[topic][:i], ps.Map[topic][i+1:]...)
				break
			}
		}
	}
}
