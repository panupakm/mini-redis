package handler

import (
	"io"

	"github.com/panupakm/miniredis/internal/payload"
	cmd "github.com/panupakm/miniredis/internal/request"
	"github.com/panupakm/miniredis/server/context"
)

func HandleSub(rw io.ReadWriter, ctx *context.Context) error {
	sub := cmd.SubReadFrom(rw)
	ps := ctx.PubSub

	ps.Sub(sub.Topic, rw)
	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err := r.WriteTo(rw)
	return err
}
