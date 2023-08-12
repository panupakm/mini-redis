//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/server"
)

func InitializeServer() *server.Server {
	wire.Build(server.NewServer, db.NewDb, pubsub.NewPubSub)
	return &server.Server{}
}
