package client

import (
	"net"

	"github.com/panupakm/miniredis/internal/payload"
)

type Subsriber struct {
	messages chan payload.SubMsg
	conn     net.Conn
}

func (s *Subsriber) NextMessage() (*payload.SubMsg, error) {
	var msg payload.SubMsg

	_, err := msg.ReadFrom(s.conn)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
