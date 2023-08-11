package handler

import (
	"io"

	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/context"
)

func HandlePub(rw io.ReadWriter, ctx *context.Context) error {
	pub := cmd.PubReadFrom(rw)
	ps := ctx.PubSub

	go ps.Pub(pub.Topic, pub.Typ, pub.Data, rw)
	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err := r.WriteTo(rw)
	return err
}
