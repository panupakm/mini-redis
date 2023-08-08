package request

import (
	"io"

	"github.com/panupakm/miniredis/payload"
)

type Ping struct {
	message string
}

const (
	PingCode = "ping"
)

func PingReadFrom(r io.Reader) *Ping {
	var msg payload.String
	msg.ReadFrom(r)
	return &Ping{
		message: string(msg),
	}
}

func (p *Ping) String() string {
	return p.message
}
