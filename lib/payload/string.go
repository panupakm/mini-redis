package payload

import (
	"encoding/binary"
	"errors"
	"io"
)

type String string

func (s String) Bytes() []byte {
	return []byte(s)
}

func (s String) String() string {
	return string(s)
}

func (s String) WriteTo(w io.Writer) (int64, error) {
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

func (s *String) ReadFrom(r io.Reader) (int64, error) {
	var typ ValueType
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

	if size == 0 {
		*s = ""
		return n, nil
	}

	buff := make([]byte, size)
	o, err := r.Read(buff)

	if err != nil {
		return n, err
	}
	*s = String(buff[:o])

	return n + int64(o), nil
}
