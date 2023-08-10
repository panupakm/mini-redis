package handler

import (
	"net"

	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/context"
)

func HandlePub(conn net.Conn, ctx *context.Context) error {
	pub := cmd.PubReadFrom(conn)
	ps := ctx.PubSub

	go ps.Pub(pub.Topic, pub.Typ, pub.Data, conn)
	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err := r.WriteTo(conn)
	return err
}
