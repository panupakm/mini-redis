package cmd

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/panupakm/miniredis/lib/payload"
)

func makeStringPayloadReader(s string) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(payload.StringType))
	binary.Write(buf, binary.BigEndian, uint32(len(s)))
	buf.WriteString(s)

	return bytes.NewReader(buf.Bytes())
}
