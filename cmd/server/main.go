// miniredis project server main.go
package main

import (
	"flag"
	"fmt"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/server"
)

func main() {
	port := flag.Uint("port", 9988, "port to listen on")
	addr := flag.String("addr", "localhost", "address to listen on")
	flag.Parse()

	s := server.NewServer(*addr, *port, db.NewDb(), pubsub.NewPubSub())
	fmt.Println("Server started")
	s.ListenAndServe()
}
