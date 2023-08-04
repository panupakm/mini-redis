package cmd

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/panupakm/miniredis/lib/payload"
)

type Pub struct {
	Topic string
	Typ   payload.ValueType
	Len   uint64
	Data  []byte
}

const (
	PubCode = "pub"
)

func PubReadFrom(conn net.Conn) *Pub {
	var topic payload.String
	topic.ReadFrom(conn)

	var typ payload.ValueType
	err := binary.Read(conn, binary.BigEndian, &typ)
	if err != nil {
		fmt.Printf("pub error: %s", err)
		return nil
	}

	var len uint32
	err = binary.Read(conn, binary.BigEndian, &len)
	if err != nil {
		fmt.Printf("pub error: %s", err)
		return nil
	}

	buff := make([]byte, len)
	_, err = conn.Read(buff)

	return &Pub{
		Topic: string(topic),
		Typ:   typ,
		Len:   uint64(len),
		Data:  buff,
	}
}

func (s *Pub) String() string {
	return fmt.Sprintf("pub topic:%s", s.Topic)
}
