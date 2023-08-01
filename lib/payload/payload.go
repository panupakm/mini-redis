package payload

import (
	"errors"
	"fmt"
	"io"
)

type ValueType uint8

const (
	BinaryType ValueType = iota + 1
	StringType
	ResultType

	MaxPayloadSize uint32 = 10 << 20
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}
