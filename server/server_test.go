// miniredis project server main.go
package server

import (
	"net"
	"testing"
	"time"

	"github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	type args struct {
		host string
		port string
		db   *db.Db
		ps   *pubsub.PubSub
	}
	db := db.NewDb()
	pubsub := pubsub.NewPubSub()
	tests := []struct {
		name string
		args args
		want Server
	}{
		{
			name: "create with localhost url",
			args: args{
				host: "localhost",
				port: "6379",
				db:   db,
				ps:   pubsub,
			},
			want: Server{
				host: "localhost",
				port: "6379",
				db:   db,
				ps:   pubsub,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.host, tt.args.port, tt.args.db, tt.args.ps)
			assert.Equal(t, tt.want, *got)
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
		{
			name: "close server with connection",
			fields: fields{
				host: "localhost",
				port: "0",
			},
			wantErr: false,
		},
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
			go s.ListenAndServe()
			time.Sleep(1 * time.Second)
			err := s.Close()
			assert.NoError(t, err)
		})
	}
}

func Test_processClient(t *testing.T) {
	type args struct {
		conn net.Conn
		ctx  *miniredis.Context
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processClient(tt.args.conn, tt.args.ctx)
		})
	}
}
