package cmd

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *Ping
	}{
		{
			name: "ping with valid message",
			args: args{
				r: makeStringPayloadReader("PONG"),
			},
			want: &Ping{
				message: "PONG",
			},
		},
		{
			name: "ping with empty message",
			args: args{
				r: makeStringPayloadReader(""),
			},
			want: &Ping{
				message: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PingReadFrom(tt.args.r)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPing_String(t *testing.T) {
	type fields struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ping with valid message",
			fields: fields{
				message: "PONG",
			},
			want: "PONG",
		},
		{
			name: "ping with empty message",
			fields: fields{
				message: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Ping{
				message: tt.fields.message,
			}
			got := p.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
