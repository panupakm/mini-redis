// miniredis project server main.go
package main

import (
	"fmt"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/server"
)

func main() {
	s := server.NewServer("localhost", "9988", db.NewDb())
	fmt.Println("Server started")
	s.Start()
}
