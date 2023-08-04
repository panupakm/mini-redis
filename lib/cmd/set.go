package cmd

import (
	"fmt"
	"net"

	"github.com/panupakm/miniredis/lib/payload"
)

type Set struct {
	Key   string
	Typ   payload.ValueType
	Value []byte
}

const (
	SetCode = "set"
)

func SetReadFrom(conn net.Conn) *Set {
	var key payload.String
	key.ReadFrom(conn)

	var value payload.String
	value.ReadFrom(conn)
	return &Set{
		Key:   key.String(),
		Value: value.Bytes(),
	}
}

func (s *Set) String() string {
	return fmt.Sprintf("key:%s value:%s", s.Key, s.Value)
}
