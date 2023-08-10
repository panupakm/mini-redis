package request

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

func SetReadFrom(r io.Reader) (Set, error) {
	var key payload.String
	key.ReadFrom(r)

	var value payload.String
	_, err := value.ReadFrom(r)
	if err != nil {
		return Set{}, err
	}
	return Set{
		Key:   key.String(),
		Typ:   payload.StringType,
		Value: value.Bytes(),
	}, nil
}

func (s *Set) String() string {
	return fmt.Sprintf("key:%s value:%s", s.Key, s.Value)
}
