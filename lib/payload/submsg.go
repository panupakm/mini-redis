package payload

import (
	"encoding/binary"
	"errors"
	"io"
)

type SubMsg struct {
	typ    ValueType
	size   uint32
	buffer []byte
}

func NewSubMsg(typ ValueType, buffer []byte) *SubMsg {
	return &SubMsg{
		typ:    typ,
		size:   uint32(len(buffer)),
		buffer: buffer,
	}
}

func (m *SubMsg) AsString() (string, bool) {
	if m.typ == StringType {
		return string(m.buffer[:m.size]), true
	}
	return "", false
}

func (m *SubMsg) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, SubMsgType)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, m.typ)
	if err != nil {
		return n, err
	}
	n++

	err = binary.Write(w, binary.BigEndian, m.size)
	if err != nil {
		return 0, err
	}
	n += 4

	o, err := w.Write(m.buffer[:m.size])
	return n + int64(o), err
}

func (m *SubMsg) ReadFrom(r io.Reader) (int64, error) {
	var typ ValueType
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != SubMsgType {
		return 0, errors.New("invalid sub message type")
	}

	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return n, err
	}
	n++

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return 0, err
	}
	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	buff := make([]byte, size)
	o, err := r.Read(buff)

	if err != nil {
		return n, err
	}

	m.size = size
	m.typ = typ
	m.buffer = buff
	return n + int64(o), nil
}
