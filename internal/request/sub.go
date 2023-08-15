package request

import (
	"bytes"
	"fmt"
	"io"

	"github.com/panupakm/miniredis/internal/payload"
)

type Sub struct {
	Topic string
}

const (
	SubCode = "sub"
)

func SubReadFrom(r io.Reader) *Sub {
	var topic payload.String
	topic.ReadFrom(r)
	return &Sub{
		Topic: string(topic),
	}
}

func (s *Sub) String() string {
	return fmt.Sprintf("topic:%s", s.Topic)
}

func (s *Sub) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	str := payload.String(SubCode)
	str.WriteTo(buf)

	str = payload.String(s.Topic)
	str.WriteTo(buf)

	return buf.Bytes()
}
