package pubsub

import (
	"github.com/panupakm/miniredis/payload"
	"net"
	"reflect"
	"testing"
)

func TestNewPubSub(t *testing.T) {
	tests := []struct {
		name string
		want *PubSub
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPubSub(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPubSub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubSub_Pub(t *testing.T) {
	type fields struct {
		Map map[string][]net.Conn
	}
	type args struct {
		topic string
		typ   payload.ValueType
		buff  []byte
		conn  net.Conn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PubSub{
				Map: tt.fields.Map,
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PubSub{
				Map: tt.fields.Map,
			}
			ps.Sub(tt.args.topic, tt.args.conn)
		})
	}
}

func TestPubSub_UnsubConnection(t *testing.T) {
	type fields struct {
		Map map[string][]net.Conn
	}
	type args struct {
		conn net.Conn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PubSub{
				Map: tt.fields.Map,
			}
			ps.UnsubConnection(tt.args.conn)
		})
	}
}
