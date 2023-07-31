package payload

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayloadSize uint32 = 10 << 20
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (b Binary) Bytes() []byte {
	return b
}

func (b Binary) String() string {
	return string(b)
}

func (b Binary) WriterTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType)
	if err != nil {
		return 0, err
	}
	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(b)))
	if err != nil {
		return 0, err
	}
	n += 4

	o, err := w.Write(b)
	return n + int64(o), err
}

func (b *Binary) ReaderFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType {
		return 0, errors.New("invalid Binary")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return 0, err
	}
	n += 4

	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	*b = make([]byte, size)
	o, err := r.Read(*b)
	return n + int64(o), err
}

type String string

func (s String) Bytes() []byte {
	return []byte(s)
}

func (s String) String() string {
	return string(s)
}

func (s String) WriterTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType)
	if err != nil {
		return 0, err
	}
	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(s)))
	if err != nil {
		return 0, err
	}
	n += 4

	o, err := w.Write([]byte(s))
	return n + int64(o), err
}

func (s *String) ReaderFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != StringType {
		return 0, errors.New("invalid string")
	}

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
	*s = String(buff[:o])

	return n + int64(o), nil
}
