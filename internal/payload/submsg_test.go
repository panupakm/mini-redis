package payload

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSubMsg(t *testing.T) {
	type args struct {
		typ    ValueType
		buffer []byte
	}
	tests := []struct {
		name string
		args args
		want SubMsg
	}{
		{
			name: "new sub message",
			args: args{
				typ:    StringType,
				buffer: []byte("hello"),
			},
			want: SubMsg{
				typ:    StringType,
				size:   5,
				buffer: []byte("hello"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSubMsg(tt.args.typ, tt.args.buffer)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestSubMsg_WriteTo(t *testing.T) {
	type fields struct {
		typ    ValueType
		size   uint32
		buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantW   string
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				typ:    StringType,
				size:   11,
				buffer: []byte("hello world"),
			},
			want:    1 + 5 + 11,
			wantW:   "\x05\x02\x00\x00\x00\x0Bhello world",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SubMsg{
				typ:    tt.fields.typ,
				size:   tt.fields.size,
				buffer: tt.fields.buffer,
			}
			w := &bytes.Buffer{}
			got, err := m.WriteTo(w)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, []byte(tt.wantW), w.Bytes())
		})
	}
}

func TestSubMsg_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantMsg string
		wantErr bool
	}{
		{
			name: "ReadFrom",
			args: args{
				r: bytes.NewReader([]byte("\x05\x02\x00\x00\x00\x0Bhello world")),
			},
			want:    17,
			wantMsg: "hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(SubMsg)
			got, err := m.ReadFrom(tt.args.r)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			gotMsg, _ := m.AsString()
			assert.Equal(t, tt.wantMsg, gotMsg)
		})
	}
}
