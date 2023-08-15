package pubsub

import (
	"io"

	"github.com/panupakm/miniredis/payload"
	"github.com/panupakm/miniredis/server/pubsub/internal"
)

type PubSub interface {
	Sub(topic string, w io.Writer)
	Pub(topic string, typ payload.ValueType, buff []byte, w io.Writer)
	Unsub(w io.Writer)
	IsSub(topic string, w io.Writer) bool
}

func NewDefaultPubSub() PubSub {
	return internal.NewDefaultPubSub()
}
