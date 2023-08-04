package handler

import (
	"fmt"
	"net"

	miniredis "github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/lib/cmd"
	"github.com/panupakm/miniredis/lib/payload"
)

type Get struct {
	key   string
	value []byte
}

func HandleGet(conn net.Conn, ctx *miniredis.Context) error {
	pair := cmd.GetReadFrom(conn)

	db := ctx.Db
	v, err := db.Get(pair.Key)
	if err != nil {
		fmt.Println("Error set:", err.Error())
		r := payload.NewErrResult(payload.StringType, []byte(err.Error()))
		_, _ = r.WriteTo(conn)
		return err
	}

	r := payload.NewResult(v.Typ, v.Buf)
	_, err = r.WriteTo(conn)
	return err
}
