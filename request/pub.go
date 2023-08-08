package request

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/panupakm/miniredis/payload"
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

func PubStringTo(w io.Writer, topic, msg string) error {
	pl := payload.String(PubCode)
	_, err := pl.WriteTo(w)
	if err != nil {
		return err
	}

	pl = payload.String(topic)
	_, err = pl.WriteTo(w)
	if err != nil {
		return err
	}

	pl = payload.String(msg)
	_, err = pl.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func PubReadFrom(r io.Reader) *Pub {
	var topic payload.String
	topic.ReadFrom(r)

	var typ payload.ValueType
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		fmt.Printf("pub error: %s", err)
		return nil
	}

	var len uint32
	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		fmt.Printf("pub error: %s", err)
		return nil
	}

	buff := make([]byte, len)
	_, err = r.Read(buff)

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