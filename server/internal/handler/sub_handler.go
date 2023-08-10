package handler

import (
	"net"

	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/context"
)

func HandleSub(conn net.Conn, ctx *context.Context) error {
	sub := cmd.SubReadFrom(conn)
	ps := ctx.PubSub

	ps.Sub(sub.Topic, conn)
	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err := r.WriteTo(conn)
	return err
}
