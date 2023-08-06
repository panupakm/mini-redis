package cmd

import (
	"fmt"
	"io"

	"github.com/panupakm/miniredis/lib/payload"
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
