// client to connect to mini redis server
package main

import (
	"fmt"
	"time"

	"github.com/panupakm/miniredis/client"
)

func main() {
	c := client.NewClient()
	_ = c.Connect(fmt.Sprintf("localhost:%s", "9988"))

	time.Sleep(2 * time.Second)
}
