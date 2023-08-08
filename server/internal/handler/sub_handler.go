package handler

import (
	"net"

	miniredis "github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
)

func HandleSub(conn net.Conn, ctx *miniredis.Context) error {
	sub := cmd.SubReadFrom(conn)
	ps := ctx.PubSub

	ps.Sub(sub.Topic, conn)
	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err := r.WriteTo(conn)
	return err
}
