package payload

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResult(t *testing.T) {
	type args struct {
		typ    ValueType
		buffer []byte
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		{
			name: "result as string",
			args: args{
				typ:    StringType,
				buffer: []byte("hello"),
			},
			want: &Result{
				Code:   0,
				Length: uint32(len("hello")),
				Typ:    StringType,
				Buffer: []byte("hello"),
			},
		},
		{
			name: "result as binary",
			args: args{
				typ:    BinaryType,
				buffer: []byte("hello"),
			},
			want: &Result{
				Code:   0,
				Length: uint32(len("hello")),
				Typ:    BinaryType,
				Buffer: []byte("hello"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewResult(tt.args.typ, tt.args.buffer)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewErrResult(t *testing.T) {
	type args struct {
		typ    ValueType
		buffer []byte
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		{
			name: "result as string",
			args: args{
				typ:    StringType,
				buffer: []byte("hello"),
			},
			want: &Result{
				Code:   1,
				Length: uint32(len("hello")),
				Typ:    StringType,
				Buffer: []byte("hello"),
			},
		},
		{
			name: "result as binary",
			args: args{
				typ:    BinaryType,
				buffer: []byte("hello"),
			},
			want: &Result{
				Code:   1,
				Length: uint32(len("hello")),
				Typ:    BinaryType,
				Buffer: []byte("hello"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrResult(tt.args.typ, tt.args.buffer)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestResult_Bytes(t *testing.T) {
	type fields struct {
		Code   uint16
		Length uint32
		Typ    ValueType
		Buffer []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "to bytes",
			fields: fields{
				Code:   0,
				Length: uint32(len("hello")),
				Typ:    StringType,
				Buffer: []byte("hello"),
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x2, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{
				Code:   tt.fields.Code,
				Length: tt.fields.Length,
				Typ:    tt.fields.Typ,
				Buffer: tt.fields.Buffer,
			}
			got := r.Bytes()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestResult_String(t *testing.T) {
	type fields struct {
		Code   uint16
		Length uint32
		Typ    ValueType
		Buffer []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "to string",
			fields: fields{
				Code:   0,
				Length: uint32(len("hello")),
				Typ:    StringType,
				Buffer: []byte("hello"),
			},
			want: fmt.Sprint("code:0 length:5 type:2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{
				Code:   tt.fields.Code,
				Length: tt.fields.Length,
				Typ:    tt.fields.Typ,
				Buffer: tt.fields.Buffer,
			}
			got := r.String()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestResult_WriteTo(t *testing.T) {
	type fields struct {
		Code   uint16
		Length uint32
		Typ    ValueType
		Buffer []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantW   string
		wantErr bool
	}{
		{
			name: "Write to",
			fields: fields{
				Code:   0,
				Length: uint32(len("hello")),
				Typ:    StringType,
				Buffer: []byte("hello"),
			},
			want:    int64(len("hello")) + 2 + 1 + 4 + 5,
			wantW:   "hello",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{
				Code:   tt.fields.Code,
				Length: tt.fields.Length,
				Typ:    tt.fields.Typ,
				Buffer: tt.fields.Buffer,
			}
			w := &bytes.Buffer{}
			got, err := r.WriteTo(w)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResult_ReadFrom(t *testing.T) {
	type fields struct {
		Code   uint16
		Length uint32
		Typ    ValueType
		Buffer []byte
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "",
			fields: fields{
				Code:   0,
				Length: uint32(len("hello")),
				Typ:    StringType,
				Buffer: []byte("hello"),
			},
			args: args{
				r: bytes.NewReader([]byte{0x03, 0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0, 0x0, 0x05, 0x02, 0x68, 0x65, 0x6c, 0x6c, 0x6f}),
			},
			want:    26,
			wantBuf: []byte("hello"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &Result{
				Code:   tt.fields.Code,
				Length: tt.fields.Length,
				Typ:    tt.fields.Typ,
				Buffer: tt.fields.Buffer,
			}
			got, err := rs.ReadFrom(tt.args.r)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantBuf, rs.Buffer)
		})
	}
}

func TestResult_DataAsString(t *testing.T) {
	type fields struct {
		Code   uint16
		Length uint32
		Typ    ValueType
		Buffer []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "non-empty string",
			fields: fields{
				Code:   0,
				Length: uint32(5),
				Typ:    StringType,
				Buffer: []byte("hello"),
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &Result{
				Code:   tt.fields.Code,
				Length: tt.fields.Length,
				Typ:    tt.fields.Typ,
				Buffer: tt.fields.Buffer,
			}
			got := rs.DataAsString()
			assert.Equal(t, tt.want, got)
		})
	}
}
