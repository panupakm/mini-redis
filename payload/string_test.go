package payload

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_Bytes(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want []byte
	}{
		{
			name: "not empty string",
			s:    "hello world",
			want: []byte("hello world"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Bytes()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestString_String(t *testing.T) {
	tests := []struct {
		name string
		s    String
		want string
	}{
		{
			name: "not empty string",
			s:    String("hello world"),
			want: "hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		s       String
		want    int64
		wantW   string
		wantErr bool
	}{
		{
			name:    "not empty string",
			s:       String("hello world"),
			want:    11 + 5,
			wantW:   "\x02\x00\x00\x00\x0Bhello world",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.s.WriteTo(w)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, []byte(tt.wantW), w.Bytes())
		})
	}
}

func TestString_ReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		s       String
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "non empty string",
			s:    String("hello world"),
			args: args{
				r: bytes.NewReader([]byte("\x02\x00\x00\x00\x0Bhello world")),
			},
			want: 11 + 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("String.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String.ReadFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
