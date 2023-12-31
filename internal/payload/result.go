package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Result struct {
	Code   uint16
	Length uint32
	Typ    ValueType
	Buffer []byte
}

func NewResult(typ ValueType, buffer []byte) *Result {
	return &Result{
		Code:   0,
		Length: uint32(len(buffer)),
		Typ:    typ,
		Buffer: buffer,
	}
}

func NewErrResult(typ ValueType, buffer []byte) *Result {
	return &Result{
		Code:   1,
		Length: uint32(len(buffer)),
		Typ:    typ,
		Buffer: buffer,
	}
}

func NewResultFromGeneral(g General) *Result {
	return &Result{
		Code:   0,
		Length: g.len,
		Typ:    g.typ,
		Buffer: g.buf,
	}
}

func (r *Result) Bytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, r.Code)
	binary.Write(&buf, binary.BigEndian, r.Length)
	binary.Write(&buf, binary.BigEndian, r.Typ)
	buf.Write(r.Buffer)
	return buf.Bytes()
}

func (r *Result) String() string {
	return fmt.Sprintf("code:%v length:%v type:%v", r.Code, r.Length, r.Typ)
}

func (r *Result) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, ResultType)
	if err != nil {
		return 0, err
	}
	var n int64 = 1
	buf := r.Bytes()
	// fmt.Println("Result Length:", len(buf), string(buf))
	err = binary.Write(w, binary.BigEndian, uint32(len(buf)))
	if err != nil {
		return 0, err
	}
	n += 4

	o, err := w.Write(buf)
	return n + int64(o), err
}

func (rs *Result) ReadFrom(r io.Reader) (int64, error) {
	var typ ValueType
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}
	rs.Typ = typ

	var n int64 = 1
	if typ != ResultType {
		return 0, errors.New("invalid result")
	}

	var size uint32
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return n, err
	}
	if err != nil {
		return 0, err
	}
	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	if err := binary.Read(r, binary.BigEndian, &rs.Code); err != nil {
		return n, err
	}
	if err := binary.Read(r, binary.BigEndian, &rs.Length); err != nil {
		return n, err
	}
	if err := binary.Read(r, binary.BigEndian, &rs.Typ); err != nil {
		return n, err
	}
	buff := make([]byte, rs.Length)
	r.Read(buff)

	rs.Buffer = buff
	if err != nil {
		return n, err
	}

	return n + int64(size), nil
}

func (rs *Result) DataAsString() string {
	if rs.Buffer == nil {
		return ""
	}
	return string(rs.Buffer)
}
