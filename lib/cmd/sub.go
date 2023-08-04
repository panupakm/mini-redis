package cmd

import (
	"fmt"
	"net"

	"github.com/panupakm/miniredis/lib/payload"
)

type Sub struct {
	Topic string
}

const (
	SubCode = "sub"
)

func SubReadFrom(conn net.Conn) *Sub {
	var topic payload.String
	topic.ReadFrom(conn)
	return &Sub{
		Topic: string(topic),
	}
}

func (s *Sub) String() string {
	return fmt.Sprintf("topic:%s", s.Topic)
}
