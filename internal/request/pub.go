package request

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/panupakm/miniredis/internal/payload"
)

type Pub struct {
	Topic string
	Typ   payload.ValueType
	Len   uint32
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

	var size uint32
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		fmt.Printf("pub error: %s", err)
		return nil
	}

	buff := make([]byte, size)
	if size > 0 {
		if _, err = r.Read(buff); err != nil {
			fmt.Printf("pub error: %s", err)
			return nil
		}
	}

	return &Pub{
		Topic: string(topic),
		Typ:   typ,
		Len:   uint32(size),
		Data:  buff,
	}
}

func (s *Pub) String() string {
	return fmt.Sprintf("pub topic:%s", s.Topic)
}

func (s *Pub) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	str := payload.String(PubCode)
	str.WriteTo(buf)
	str = payload.String(s.Topic)
	str.WriteTo(buf)

	_ = binary.Write(buf, binary.BigEndian, s.Typ)
	_ = binary.Write(buf, binary.BigEndian, s.Len)
	buf.Write(s.Data)

	return buf.Bytes()
}
