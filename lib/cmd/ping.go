package cmd

import (
	"net"

	"github.com/panupakm/miniredis/lib/payload"
)

type Ping struct {
	message string
}

const (
	PingCode = "ping"
)

func PingReadFrom(conn net.Conn) *Ping {
	var msg payload.String
	msg.ReadFrom(conn)
	return &Ping{
		message: msg.String(),
	}
}

func (p *Ping) String() string {
	return p.message
}
