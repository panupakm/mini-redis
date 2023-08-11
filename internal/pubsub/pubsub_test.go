package pubsub

import (
	"net"
	"testing"

	"github.com/panupakm/miniredis/mock"
	"github.com/panupakm/miniredis/payload"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewPubSub(t *testing.T) {
	tests := []struct {
		name string
		want *PubSub
	}{
		{
			name: "new pubsub",
			want: &PubSub{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPubSub()
			assert.IsType(t, &PubSub{}, got)
		})
	}
}

func TestPubSub_Pub(t *testing.T) {
	type args struct {
		topic string
		typ   payload.ValueType
		buff  []byte
		conn  net.Conn
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "pub with topic and string value",
			args: args{
				topic: "topic",
				typ:   payload.StringType,
				buff:  []byte("test"),
				conn:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mconn := mock.NewMockConn(ctrl)
			mconn.EXPECT().Write(gomock.Any()).AnyTimes()
			connmap := map[string][]net.Conn{
				"topic": {mconn},
			}
			ps := &PubSub{
				connmap: connmap,
			}
			ps.Pub(tt.args.topic, tt.args.typ, tt.args.buff, tt.args.conn)
		})
	}
}

func TestPubSub_Sub(t *testing.T) {
	type fields struct {
		Map map[string][]net.Conn
	}
	type args struct {
		topic string
		conn  net.Conn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "sub with topic and string value",
			fields: fields{
				Map: make(map[string][]net.Conn),
			},
			args: args{
				topic: "topic",
				conn:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			connMock := mock.NewMockConn(ctrl)
			ps := &PubSub{
				connmap: tt.fields.Map,
			}
			ps.Sub(tt.args.topic, connMock)
			assert.True(t, ps.isSub(tt.args.topic))
		})
	}
}

func TestPubSub_Unsub(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "unsub with valid connection",
			args: args{
				topic: "topic",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockconn := mock.NewMockConn(ctrl)

			ps := &PubSub{
				connmap: map[string][]net.Conn{
					"topic": {mockconn},
				},
			}
			assert.True(t, ps.isSub("topic"))
			ps.Unsub(mockconn)
			assert.False(t, ps.isSub("topic"))
		})
	}
}
