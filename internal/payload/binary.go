package payload

import (
	"encoding/binary"
	"errors"
	"io"
)

type Binary []byte

func (b *Binary) Bytes() []byte {
	return *b
}

func (b *Binary) String() string {
	return string(*b)
}

func (b *Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType)
	if err != nil {
		return 0, err
	}
	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(*b)))
	if err != nil {
		return 0, err
	}
	n += 4

	o, err := w.Write(*b)
	return n + int64(o), err
}

func (b *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ ValueType
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
