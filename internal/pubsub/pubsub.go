package pubsub

import (
	"fmt"
	"io"
	"sync"

	"github.com/panupakm/miniredis/payload"
)

type PubSub struct {
	writermap map[string][]io.Writer
	mu        sync.RWMutex
}

type PubSuber interface {
	Sub(topic string, w io.Writer)
	Pub(topic string, typ payload.ValueType, buff []byte, w io.Writer)
	Unsub(w io.Writer)
	IsSub(topic string, w io.Writer) bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		writermap: make(map[string][]io.Writer),
	}
}

func (ps *PubSub) Sub(topic string, w io.Writer) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.writermap[topic] = append(ps.writermap[topic], w)
}

func (ps *PubSub) isSub(topic string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	conns, ok := ps.writermap[topic]
	return ok && len(conns) > 0
}

func (ps *PubSub) Pub(topic string, typ payload.ValueType, buff []byte, w io.Writer) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	conns := ps.writermap[topic]
	for _, con := range conns {
		if con == w {
			continue
		}

		msg := payload.NewSubMsg(typ, buff)
		_, err := msg.WriteTo(con)
		if err != nil {
			fmt.Println("pub error:", err)
		}
	}
}

func (ps *PubSub) Unsub(conn io.Writer) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for topic, conns := range ps.writermap {
		for i, c := range conns {
			if c == conn {
				ps.writermap[topic] = append(ps.writermap[topic][:i], ps.writermap[topic][i+1:]...)
				break
			}
		}
	}
}
