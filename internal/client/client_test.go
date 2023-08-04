// client to connect to mini redis server
package client

import (
	"net"
	"reflect"
	"testing"

	"github.com/panupakm/miniredis/lib/payload"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Connect(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				conn: tt.fields.conn,
				addr: tt.fields.addr,
			}
			if err := c.Connect(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("Client.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				conn: tt.fields.conn,
				addr: tt.fields.addr,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Ping(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    chan ResultChannel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				conn: tt.fields.conn,
				addr: tt.fields.addr,
			}
			got, err := c.Ping(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SetString(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    chan ResultChannel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				conn: tt.fields.conn,
				addr: tt.fields.addr,
			}
			got, err := c.SetString(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    chan ResultChannel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				conn: tt.fields.conn,
				addr: tt.fields.addr,
			}
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ReadResult(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *payload.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				conn: tt.fields.conn,
				addr: tt.fields.addr,
			}
			got, err := c.ReadResult()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ReadResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ReadResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
