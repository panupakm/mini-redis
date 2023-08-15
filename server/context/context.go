package context

import (
	"context"

	"github.com/panupakm/miniredis/server/pubsub"
	"github.com/panupakm/miniredis/server/storage"
)

type Context struct {
	context.Context
	Storage storage.Storage
	PubSub  *pubsub.PubSub
}

func NewContext(storage storage.Storage, ps *pubsub.PubSub) *Context {
	return &Context{
		Context: context.Background(),
		Storage: storage,
		PubSub:  ps,
	}
}
