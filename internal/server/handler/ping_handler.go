package handler

import (
	"net"

	"github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/lib/cmd"
	"github.com/panupakm/miniredis/lib/payload"
)

func HandlePing(conn net.Conn, ctx *miniredis.Context) error {
	ping := cmd.PingReadFrom(conn)
	msg := ping.String()
	result := payload.NewResult(payload.StringType, []byte(msg))

	_, err := result.WriteTo(conn)
	return err
}
