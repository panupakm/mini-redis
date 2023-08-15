package handler

import (
	"fmt"
	"io"

	"github.com/panupakm/miniredis/internal/payload"
	cmd "github.com/panupakm/miniredis/internal/request"
	"github.com/panupakm/miniredis/server/context"
)

func HandleSet(rw io.ReadWriter, ctx *context.Context) error {
	pair := cmd.SetReadFrom(rw)
	storage := ctx.Storage
	err := storage.Set(pair.Key, *payload.NewGeneral(payload.StringType, pair.Value))
	if err != nil {
		fmt.Println("Error set:", err.Error())
		payload.NewErrResult(payload.StringType, []byte(err.Error()))
		return err
	}

	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err = r.WriteTo(rw)
	return err
}
