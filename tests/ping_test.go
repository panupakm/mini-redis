package tests

import (
	"testing"

	"github.com/panupakm/mini-redis/client"
	"github.com/panupakm/mini-redis/server"
)

func TestServer(t *testing.T) {
	t.Log("Start server...")
	go server.Start("tcp", "localhost", "9988")

	t.Log("Start client...")
	client.Start("tcp", "localhost", "9988")
}
