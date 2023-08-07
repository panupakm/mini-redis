package cmd

import (
	"io"
	"testing"

	"github.com/panupakm/miniredis/payload"
	"github.com/stretchr/testify/assert"
)

func TestGetReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *Get
	}{
		{
			name: "get with valid key",
			args: args{
				r: payload.MakeStringPayloadReader("Yahoo"),
			},
			want: &Get{
				Key: "Yahoo",
			},
		},
		{
			name: "get with empty key",
			args: args{
				r: payload.MakeStringPayloadReader(""),
			},
			want: &Get{
				Key: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get := GetReadFrom(tt.args.r)
			if tt.want == nil {

			}
			assert.Equal(t, tt.want, get)
		})
	}
}

func TestGet_String(t *testing.T) {
	type fields struct {
		Key string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "valid key",
			fields: fields{
				Key: "Yahoo",
			},
			want: "key:Yahoo",
		},
		{
			name: "empty key",
			fields: fields{
				Key: "",
			},
			want: "key:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Get{
				Key: tt.fields.Key,
			}
			assert.Equal(t, tt.want, g.String())
		})
	}
}
