package handler

import (
	"fmt"
	"net"

	miniredis "github.com/panupakm/miniredis"
	mdb "github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
)

type Set struct {
	key   string
	value []byte
}

func HandleSet(conn net.Conn, ctx *miniredis.Context) error {
	pair := cmd.SetReadFrom(conn)

	db := ctx.Db
	err := db.Set(pair.Key, mdb.NewTypeBuffer(payload.StringType, pair.Value))
	if err != nil {
		fmt.Println("Error set:", err.Error())
		payload.NewErrResult(payload.StringType, []byte(err.Error()))
		return err
	}

	r := payload.NewResult(payload.StringType, []byte("OK"))
	_, err = r.WriteTo(conn)
	return err
}
