package ping

import (
	"fmt"
	"net"

	"github.com/panupakm/mini-redis/lib/payload"
)

type Ping struct {
	message string
}

const (
	Code = "ping"
)

func NewPing(conn net.Conn) *Ping {
	var msg payload.String
	msg.ReaderFrom(conn)
	fmt.Println("Message: ", msg)
	return &Ping{
		message: msg.String(),
	}
}

func (p *Ping) String() string {
	return p.message
}
