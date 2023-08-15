package payload

import (
	"bytes"
	"encoding/binary"
	"io"
)

func MakeStringPayloadReader(s string) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(StringType))
	binary.Write(buf, binary.BigEndian, uint32(len(s)))
	buf.WriteString(s)

	return bytes.NewReader(buf.Bytes())
}

func MakeBinaryPayloadReader(b []byte) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(BinaryType))
	binary.Write(buf, binary.BigEndian, uint32(len(b)))
	buf.Write(b)

	return bytes.NewReader(buf.Bytes())
}
