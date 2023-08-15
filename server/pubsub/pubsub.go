package pubsub

import (
	"io"

	"github.com/panupakm/miniredis/internal/payload"
	"github.com/panupakm/miniredis/server/pubsub/internal"
)

type PubSub interface {
	Sub(topic string, w io.Writer)
	Pub(topic string, typ payload.ValueType, buff []byte, w io.Writer)
	Unsub(w io.Writer)
	IsSub(topic string) bool
}

func NewDefaultPubSub() PubSub {
	return internal.NewPubSub()
}
