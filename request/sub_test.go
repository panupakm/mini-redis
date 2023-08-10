package request

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *Sub
	}{
		{
			name: "Sub to non empty topic",
			args: args{
				r: bytes.NewReader([]byte("\x02\x00\x00\x00\x08greeting")),
			},
			want: &Sub{
				Topic: "greeting",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SubReadFrom(tt.args.r)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSub_String(t *testing.T) {
	type fields struct {
		Topic string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Sub to non empty topic",
			fields: fields{
				Topic: "greeting",
			},
			want: "topic:greeting",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sub{
				Topic: tt.fields.Topic,
			}
			got := s.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
