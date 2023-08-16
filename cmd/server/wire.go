//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/panupakm/miniredis/server"
	"github.com/panupakm/miniredis/server/pubsub"
	"github.com/panupakm/miniredis/server/storage"
)

func InitializeServer() *server.Server {
	wire.Build(server.NewServer, storage.NewDefaultStorage, pubsub.NewDefaultPubSub)
	return &server.Server{}
}
