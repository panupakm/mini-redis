package request

import (
	"bytes"
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
		Typ:   payload.StringType,
		Value: value.Bytes(),
	}
}

func (s *Set) String() string {
	return fmt.Sprintf("key:%s value:%s", s.Key, s.Value)
}

func (s *Set) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	str := payload.String(SetCode)
	str.WriteTo(buf)

	str = payload.String(s.Key)
	str.WriteTo(buf)

	str = payload.String(s.Value)
	str.WriteTo(buf)

	return buf.Bytes()
}
