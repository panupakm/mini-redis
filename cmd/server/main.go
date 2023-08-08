// miniredis project server main.go
package main

import (
	"fmt"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/server"
)

func main() {
	s := server.NewServer("localhost", "9988", db.NewDb(), pubsub.NewPubSub())
	fmt.Println("Server started")
	s.ListenAndServe()
}
