package handler

import (
	"fmt"
	"io"

	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/context"
)

type Set struct {
	key   string
	value []byte
}

func HandleSet(rw io.ReadWriter, ctx *context.Context) error {
	pair := cmd.SetReadFrom(rw)
	db := ctx.Db
	err := db.Set(pair.Key, *payload.NewGeneral(payload.StringType, pair.Value))
	if err != nil {
		fmt.Println("Error set:", err.Error())
		payload.NewErrResult(payload.StringType, []byte(err.Error()))
		return err
	}

	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err = r.WriteTo(rw)
	return err
}
