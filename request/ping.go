package request

import (
	"bytes"
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

func (p *Ping) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	str := payload.String(PingCode)
	str.WriteTo(buf)
	str = payload.String(p.message)
	str.WriteTo(buf)
	return buf.Bytes()
}
