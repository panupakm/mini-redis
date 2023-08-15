package pubsub

import (
	"io"
	"net"
	"testing"

	"github.com/panupakm/miniredis/internal/mock"
	"github.com/panupakm/miniredis/internal/payload"
	"github.com/panupakm/miniredis/server/pubsub/internal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewPubSub(t *testing.T) {
	tests := []struct {
		name string
		want PubSub
	}{
		{
			name: "new pubsub",
			want: NewDefaultPubSub(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDefaultPubSub()
			assert.IsType(t, NewDefaultPubSub(), got)
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
			writermap := map[string][]io.Writer{
				"topic": {mconn},
			}
			ps := internal.NewPubSubWithWriterMap(writermap)
			ps.Pub(tt.args.topic, tt.args.typ, tt.args.buff, tt.args.conn)
		})
	}
}

func TestPubSub_Sub(t *testing.T) {
	type fields struct {
		Map map[string][]io.Writer
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
				Map: make(map[string][]io.Writer),
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
			ps := internal.NewPubSubWithWriterMap(tt.fields.Map)
			ps.Sub(tt.args.topic, connMock)
			assert.True(t, ps.IsSub(tt.args.topic))
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

			ps := internal.NewPubSubWithWriterMap(map[string][]io.Writer{
				"topic": {mockconn},
			})
			assert.True(t, ps.IsSub("topic"))
			ps.Unsub(mockconn)
			assert.False(t, ps.IsSub("topic"))
		})
	}
}
