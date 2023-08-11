package request

import (
	"bytes"
	"fmt"
	"io"

	"github.com/panupakm/miniredis/payload"
)

type Get struct {
	Key string
}

const (
	GetCode = "get"
)

func GetReadFrom(r io.Reader) *Get {
	var key payload.String
	key.ReadFrom(r)
	return &Get{
		Key: key.String(),
	}
}

func (g *Get) String() string {
	return fmt.Sprintf("key:%s", g.Key)
}

func (g *Get) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	str := payload.String(GetCode)
	str.WriteTo(buf)

	str = payload.String(g.Key)
	str.WriteTo(buf)

	return buf.Bytes()
}
