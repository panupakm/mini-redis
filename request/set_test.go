package request

import (
	"bytes"
	"io"
	"testing"

	"github.com/panupakm/miniredis/payload"
	"github.com/stretchr/testify/assert"
)

func TestSetReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantSet Set
	}{
		{
			name: "string key and string value",
			args: args{
				r: bytes.NewBuffer([]byte("\x02\x00\x00\x00\x08greeting\x02\x00\x00\x00\x0Bhello world")),
			},
			wantSet: Set{
				Key:   "greeting",
				Typ:   payload.StringType,
				Value: []byte("hello world"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetReadFrom(tt.args.r)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSet, got)
		})
	}
}

func TestSet_String(t *testing.T) {
	type fields struct {
		Key   string
		Typ   payload.ValueType
		Value []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string key and string value",
			fields: fields{
				Key:   "greeting",
				Typ:   payload.StringType,
				Value: []byte("hello world"),
			},
			want: "key:greeting value:hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Set{
				Key:   tt.fields.Key,
				Typ:   tt.fields.Typ,
				Value: tt.fields.Value,
			}
			got := s.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
