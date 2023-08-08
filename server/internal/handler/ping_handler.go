package handler

import (
	"net"

	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
)

func HandlePing(conn net.Conn) error {
	ping := cmd.PingReadFrom(conn)
	msg := ping.String()
	if msg == "" {
		msg = "pong"
	}
	result := payload.NewResult(payload.StringType, []byte(msg))

	_, err := result.WriteTo(conn)
	return err
}
