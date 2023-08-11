package handler

import (
	"io"

	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
)

func HandlePing(rw io.ReadWriter) error {
	ping := cmd.PingReadFrom(rw)
	msg := ping.String()
	if msg == "" {
		msg = "pong"
	}
	result := payload.NewResult(payload.StringType, []byte(msg))

	_, err := result.WriteTo(rw)
	return err
}
