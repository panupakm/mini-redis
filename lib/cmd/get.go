package cmd

import (
	"fmt"
	"net"

	"github.com/panupakm/miniredis/lib/payload"
)

type Get struct {
	Key string
}

const (
	GetCode = "get"
)

func GetReadFrom(conn net.Conn) *Get {
	var key payload.String
	key.ReadFrom(conn)
	return &Get{
		Key: key.String(),
	}
}

func (g *Get) String() string {
	return fmt.Sprintf("key:%s", g.Key)
}
