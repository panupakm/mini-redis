package internal

import (
	"fmt"
	"io"
	"sync"

	"github.com/panupakm/miniredis/internal/payload"
)

type DefaultPubSub struct {
	writermap map[string][]io.Writer
	mu        sync.RWMutex
}

func NewPubSub() *DefaultPubSub {
	return &DefaultPubSub{
		writermap: make(map[string][]io.Writer),
	}
}

func NewPubSubWithWriterMap(writermap map[string][]io.Writer) *DefaultPubSub {
	return &DefaultPubSub{
		writermap: writermap,
	}
}

func (ps *DefaultPubSub) Sub(topic string, w io.Writer) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.writermap[topic] = append(ps.writermap[topic], w)
}

func (ps *DefaultPubSub) IsSub(topic string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	conns, ok := ps.writermap[topic]
	return ok && len(conns) > 0
}

func (ps *DefaultPubSub) Pub(topic string, typ payload.ValueType, buff []byte, w io.Writer) {
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

func (ps *DefaultPubSub) Unsub(conn io.Writer) {
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
