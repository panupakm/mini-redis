package cmd

import (
	"fmt"
	"net"

	"github.com/panupakm/miniredis/lib/payload"
)

type Get struct {
	Key string
}

const (
	GetCode = "get"
)

func GetWriteTo(key string, conn net.Conn) (int64, error) {
	pl := payload.String(GetCode)
	n, err := pl.WriterTo(conn)
	if err != nil {
		return n, err
	}

	pl = payload.String(key)
	o, err := pl.WriterTo(conn)
	return int64(o) + n, err
}

func GetReadFrom(conn net.Conn) *Get {
	var key payload.String
	key.ReaderFrom(conn)
	return &Get{
		Key: key.String(),
	}
}

func (g *Get) String() string {
	return fmt.Sprintf("key:%s", g.Key)
}
