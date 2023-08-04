package miniredis

import (
	"context"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
)

type Context struct {
	context.Context
	Db     *db.Db
	PubSub *pubsub.PubSub
}

func NewContext(db *db.Db, ps *pubsub.PubSub) *Context {
	return &Context{
		Context: context.Background(),
		Db:      db,
		PubSub:  ps,
	}
}
