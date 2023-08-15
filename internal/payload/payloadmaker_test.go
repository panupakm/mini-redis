package payload

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeStringPayloadReader(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want io.Reader
	}{
		{
			name: "non-empty string payload",
			args: args{
				s: "test",
			},
			want: io.NopCloser(strings.NewReader("\x02\x00\x00\x00\x04test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeStringPayloadReader(tt.args.s)
			w := bytes.Buffer{}
			tt.want.Read(w.Bytes())

			g := bytes.Buffer{}
			got.Read(g.Bytes())
			assert.Equal(t, w, g)
		})
	}
}

func TestMakeBinaryPayloadReader(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want io.Reader
	}{
		{
			name: "non-empty binary payload",
			args: args{
				b: []byte{3, 1, 2, 3},
			},
			want: io.NopCloser(bytes.NewReader([]byte{1, 0, 0, 0, 3, 1, 2, 3})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeBinaryPayloadReader(tt.args.b)
			w := bytes.Buffer{}
			tt.want.Read(w.Bytes())

			g := bytes.Buffer{}
			got.Read(g.Bytes())
			assert.Equal(t, w, g)
		})
	}
}
