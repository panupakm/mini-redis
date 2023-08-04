package common

import (
	"fmt"
	"testing"

	"github.com/panupakm/miniredis/internal/client"
	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/internal/server"
	"github.com/stretchr/testify/require"
)

func SetUpServer(t *testing.T, port string) *server.Server {

	d := db.NewDb()
	ps := pubsub.NewPubSub()

	t.Log("Start server...")
	s := server.NewServer("localhost", port, d, ps)
	go s.Start()

	return s
}

func SetUpClient(t *testing.T, port string) (*client.Client, func()) {
	c := client.NewClient()
	err := c.Connect(fmt.Sprintf("localhost:%s", port))
	require.NoError(t, err)
	return c, func() {
		c.Close()
	}
}

func SetUpServerClient(t *testing.T) (*server.Server, *client.Client, func()) {

	d := db.NewDb()
	ps := pubsub.NewPubSub()

	t.Log("Start server...")
	s := server.NewServer("localhost", "9988", d, ps)
	go s.Start()

	c := client.NewClient()
	err := c.Connect("localhost:9988")
	require.NoError(t, err)
	return s, c, func() {
		c.Close()
	}
}
