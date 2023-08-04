package handler

import (
	"net"

	miniredis "github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/lib/cmd"
	"github.com/panupakm/miniredis/lib/payload"
)

func HandlePub(conn net.Conn, ctx *miniredis.Context) error {
	pub := cmd.PubReadFrom(conn)
	ps := ctx.PubSub

	go ps.Pub(pub.Topic, pub.Typ, pub.Data, conn)
	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err := r.WriteTo(conn)
	return err
}
