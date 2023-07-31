package common

import (
	"testing"

	"github.com/panupakm/mini-redis/client"
	"github.com/panupakm/mini-redis/server"
	"github.com/stretchr/testify/require"
)

func SetUpServerClient(t *testing.T) (*server.Server, *client.Client, func()) {
	t.Log("Start server...")
	s := server.NewServer("localhost", "9988")
	go s.Start()

	c := client.NewClient()
	err := c.Connect("localhost:9988")
	require.NoError(t, err)
	return s, c, func() {
		c.Close()
	}
}
