// mini-redis project server main.go
package main

import (
	"fmt"

	"github.com/panupakm/mini-redis/server"
)

func main() {
	server.NewServer("localhost", "9988")
	fmt.Println("Server started")
}
