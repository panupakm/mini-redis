package cmd

import (
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
