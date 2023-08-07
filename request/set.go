package cmd

import (
	"fmt"
	"io"

	"github.com/panupakm/miniredis/payload"
)

type Set struct {
	Key   string
	Typ   payload.ValueType
	Value []byte
}

const (
	SetCode = "set"
)

func SetReadFrom(r io.Reader) *Set {
	var key payload.String
	key.ReadFrom(r)

	var value payload.String
	value.ReadFrom(r)
	return &Set{
		Key:   key.String(),
		Value: value.Bytes(),
	}
}

func (s *Set) String() string {
	return fmt.Sprintf("key:%s value:%s", s.Key, s.Value)
}
