package miniredis

import (
	"context"

	"github.com/panupakm/miniredis/internal/db"
)

type Context struct {
	context.Context
	Db *db.Db
}

func NewContext(db *db.Db) *Context {
	return &Context{
		Context: context.Background(),
		Db:      db,
	}
}
