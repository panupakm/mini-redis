package common

import (
	"fmt"
	"testing"

	"github.com/panupakm/miniredis/client"
	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/server"
	"github.com/stretchr/testify/require"
)

func SetUpServer(t *testing.T, port uint, config *server.Config) *server.Server {

	d := db.NewDb()
	ps := pubsub.NewPubSub()

	t.Log("Start server...")
	s := server.NewServer(d, ps)
	go s.ListenAndServe("localhost", port, config)

	return s
}

func SetUpClient(t *testing.T, port uint) *client.Client {
	c := client.NewClient()
	err := c.Connect(fmt.Sprintf("localhost:%d", port), nil)
	require.NoError(t, err)
	return c
}

func SetUpServerClient(t *testing.T) (*server.Server, *client.Client, func()) {

	d := db.NewDb()
	ps := pubsub.NewPubSub()

	t.Log("Start server...")
	s := server.NewServer(d, ps)
	go s.ListenAndServe("localhost", 9988, nil)

	c := client.NewClient()
	err := c.Connect("localhost:9988", nil)
	require.NoError(t, err)
	return s, c, func() {
		c.Close()
	}
}
