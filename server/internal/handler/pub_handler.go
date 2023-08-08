package handler

import (
	"net"

	miniredis "github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
)

func HandlePub(conn net.Conn, ctx *miniredis.Context) error {
	pub := cmd.PubReadFrom(conn)
	ps := ctx.PubSub

	go ps.Pub(pub.Topic, pub.Typ, pub.Data, conn)
	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err := r.WriteTo(conn)
	return err
}
