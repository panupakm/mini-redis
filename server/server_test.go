// miniredis project server main.go
package server

import (
	"net"
	"reflect"
	"testing"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
)

func TestNewServer(t *testing.T) {
	type args struct {
		host string
		port string
		db   *db.Db
		ps   *pubsub.PubSub
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.host, tt.args.port, tt.args.db, tt.args.ps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Close(t *testing.T) {
	type fields struct {
		host     string
		port     string
		conn     net.Conn
		listener net.Listener
		db       *db.Db
		ps       *pubsub.PubSub
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
			s := &Server{
				host:     tt.fields.host,
				port:     tt.fields.port,
				conn:     tt.fields.conn,
				listener: tt.fields.listener,
				db:       tt.fields.db,
				ps:       tt.fields.ps,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	type fields struct {
		host     string
		port     string
		conn     net.Conn
		listener net.Listener
		db       *db.Db
		ps       *pubsub.PubSub
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				host:     tt.fields.host,
				port:     tt.fields.port,
				conn:     tt.fields.conn,
				listener: tt.fields.listener,
				db:       tt.fields.db,
				ps:       tt.fields.ps,
			}
			s.ListenAndServe()
		})
	}
}

// func TestServer_processClient(t *testing.T) {
// 	type fields struct {
// 		host     string
// 		port     string
// 		conn     net.Conn
// 		listener net.Listener
// 		db       *db.Db
// 		ps       *pubsub.PubSub
// 	}
// 	type args struct {
// 		conn net.Conn
// 		ctx  *miniredis.Context
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Server{
// 				host:     tt.fields.host,
// 				port:     tt.fields.port,
// 				conn:     tt.fields.conn,
// 				listener: tt.fields.listener,
// 				db:       tt.fields.db,
// 				ps:       tt.fields.ps,
// 			}
// 			s.processClient(tt.args.conn, tt.args.ctx)
// 		})
// 	}
// }
