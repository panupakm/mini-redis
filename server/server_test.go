// miniredis project server main.go
package server

import (
	"crypto/tls"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/mock"
	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	scontext "github.com/panupakm/miniredis/server/context"
	"github.com/panupakm/miniredis/server/internal/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewServer(t *testing.T) {
	type args struct {
		host string
		port uint
		db   *db.Db
		ps   *pubsub.PubSub
	}
	db := db.NewDb()
	pubsub := pubsub.NewPubSub()
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "create with localhost url",
			args: args{
				host: "localhost",
				port: 6379,
				db:   db,
				ps:   pubsub,
			},
			want: &Server{
				host: "localhost",
				port: "6379",
				db:   db,
				ps:   pubsub,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.host, tt.args.port, tt.args.db, tt.args.ps, &tls.Config{})
			assert.Equal(t, tt.want.host, got.host)
			assert.Equal(t, tt.want.port, got.port)
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
				host:      tt.fields.host,
				port:      tt.fields.port,
				conn:      tt.fields.conn,
				listener:  tt.fields.listener,
				db:        tt.fields.db,
				ps:        tt.fields.ps,
				closechan: make(chan struct{}),
			}
			go s.ListenAndServe()
			time.Sleep(1 * time.Second)
			err := s.Close()
			assert.NoError(t, err)
		})
	}
}

func testHandlerHelper(ctrl *gomock.Controller, ctx *scontext.Context, code string) {
	mockHandler := handler.NewMockHandler(ctrl)
	mockConn := mock.NewMockConn(ctrl)
	mockAddr := mock.NewMockAddr(ctrl)

	switch code {
	case cmd.GetCode:
		mockHandler.EXPECT().HandleGet(mockConn, ctx).Times(1).Return(nil)
	case cmd.SetCode:
		mockHandler.EXPECT().HandleSet(mockConn, ctx).Times(1).Return(nil)
	case cmd.PingCode:
		mockHandler.EXPECT().HandlePing(mockConn).Times(1).Return(nil)
	case cmd.SubCode:
		mockHandler.EXPECT().HandleSub(mockConn, ctx).Times(1).Return(nil)
	case cmd.PubCode:
		mockHandler.EXPECT().HandlePub(mockConn, ctx).Times(1).Return(nil)
	}

	callCount := 0
	mockConn.EXPECT().RemoteAddr().Return(mockAddr).AnyTimes()
	mockAddr.EXPECT().String().Return("127.0.0.1:23").Times(1)
	mockConn.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (n int, err error) {
		callCount++
		switch callCount {
		case 1:
			p[0] = byte(payload.StringType)
			return 1, nil
		case 2:
			binary.BigEndian.PutUint32(p, uint32(len(code)))
			return 4, nil
		case 3:
			copy(p, []byte(code))
			return len(code), nil
		default:
			return 0, io.EOF
		}
	}).AnyTimes()
	processClient(mockConn, ctx, mockHandler, nil)
}

func Test_processClientHandle(t *testing.T) {
	type args struct {
		ctx  *scontext.Context
		code string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handle get success",
			args: args{
				ctx:  &scontext.Context{},
				code: cmd.GetCode,
			},
		},
		{
			name: "handle ping success",
			args: args{
				ctx:  &scontext.Context{},
				code: cmd.PingCode,
			},
		},
		{
			name: "handle pub success",
			args: args{
				ctx:  &scontext.Context{},
				code: cmd.PubCode,
			},
		},
		{
			name: "handle sub success",
			args: args{
				ctx:  &scontext.Context{},
				code: cmd.SubCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			testHandlerHelper(ctrl, tt.args.ctx, tt.args.code)
		})
	}
}
