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
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.b.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Binary.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Binary.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Binary.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
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
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Binary.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Binary.ReadFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
