// Package handler is internal. It implements how to handle multiple requests
package handler

import (
	"fmt"
	"io"

	"github.com/panupakm/miniredis/internal/payload"
	cmd "github.com/panupakm/miniredis/internal/request"
	"github.com/panupakm/miniredis/server/context"
)

func HandleGet(rw io.ReadWriter, ctx *context.Context) error {
	pair := cmd.GetReadFrom(rw)

	storage := ctx.Storage
	v, err := storage.Get(pair.Key)
	if err != nil {
		fmt.Println("Error set:", err.Error())
		r := payload.NewErrResult(payload.StringType, []byte(err.Error()))
		_, _ = r.WriteTo(rw)
		return err
	}

	r := payload.NewResultFromGeneral(v)
	_, err = r.WriteTo(rw)
	return err
}
