package payload

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinary_Bytes(t *testing.T) {
	tests := []struct {
		name string
		b    Binary
		want []byte
	}{
		{
			name: "modify binary payload",
			b:    Binary{0x01, 0x02, 0x03},
			want: []byte{0x01, 0x02, 0x04},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Bytes()[2] = 0x04
			assert.Equal(t, tt.b.Bytes(), tt.want)
		})
	}
}

func TestBinary_String(t *testing.T) {
	tests := []struct {
		name string
		b    Binary
		want string
	}{
		{
			name: "valid string",
			b:    Binary("Hello World!"),
			want: "Hello World!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.b.String(), tt.want)
		})
	}
}

func TestBinary_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		b       Binary
		want    int64
		wantW   string
		wantErr bool
	}{
		{
			name:    "valid string",
			b:       Binary("Hello World!"),
			want:    12 + 5,
			wantW:   "\x01\x00\x00\x00\fHello World!",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.b.WriteTo(w)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			gotW := w.String()
			assert.Equal(t, tt.wantW, gotW)
		})
	}
}

func TestBinary_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		b       *Binary
		args    args
		wantGot int64
		wantB   []byte
		wantErr bool
	}{
		{
			name: "read from",
			b:    &Binary{},
			args: args{
				r: MakeBinaryPayloadReader([]byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}),
			},
			wantGot: 12,
			wantB:   []byte{0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
		{
			name: "read from over maximum length",
			b:    &Binary{},
			args: args{
				r: MakeBinaryPayloadReader(make([]byte, MaxPayloadSize+1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.ReadFrom(tt.args.r)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantGot, got)
				assert.Equal(t, tt.wantB, tt.b.Bytes())
			}
		})
	}
}
