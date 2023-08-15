// miniredis project server main.go
package server

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/panupakm/miniredis/mock"
	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/pubsub"

	// "github.com/panupakm/miniredis/server/context"
	scontext "github.com/panupakm/miniredis/server/context"
	"github.com/panupakm/miniredis/server/internal/handler"
	"github.com/panupakm/miniredis/server/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewServer(t *testing.T) {
	type args struct {
		host    string
		port    uint
		storage storage.Storage
		pubsub  *pubsub.PubSub
	}
	storage := storage.NewDefaultStorage()
	pubsub := pubsub.NewPubSub()
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "create with localhost url",
			args: args{
				storage: storage,
				pubsub:  pubsub,
			},
			want: &Server{
				storage: storage,
				pubsub:  pubsub,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.storage, tt.args.pubsub)
			assert.Equal(t, tt.want.storage, got.storage)
			assert.Equal(t, tt.want.pubsub, got.pubsub)
		})
	}
}

func TestServer_Close(t *testing.T) {
	type fields struct {
		host     string
		port     uint
		conn     net.Conn
		listener net.Listener
		storage  storage.Storage
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
				port: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				conn:      tt.fields.conn,
				listener:  tt.fields.listener,
				storage:   tt.fields.storage,
				pubsub:    tt.fields.ps,
				closechan: make(chan struct{}),
			}
			go s.ListenAndServe(tt.fields.host, tt.fields.port, nil)
			time.Sleep(1 * time.Second)
			err := s.Close()
			assert.NoError(t, err)
		})
	}
}

type bytesQueue struct {
	queue [][]byte
}

func (b *bytesQueue) Add(data []byte) *bytesQueue {
	b.queue = append(b.queue, data[:1])
	b.queue = append(b.queue, data[1:5])
	b.queue = append(b.queue, data[5:])
	return b
}

func (b *bytesQueue) Next() (data []byte) {
	if len(b.queue) > 0 {
		data = b.queue[0]
		b.queue = b.queue[1:]
		return data
	}
	return nil
}

func testHandlerHelper(t *testing.T, ctx *scontext.Context, code string, q *bytesQueue) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := handler.NewMockHandler(ctrl)
	mockConn := mock.NewMockConn(ctrl)
	mockAddr := mock.NewMockAddr(ctrl)

	mockConn.EXPECT().RemoteAddr().Return(mockAddr).AnyTimes()
	mockAddr.EXPECT().String().Return("127.0.0.1:23").AnyTimes()
	mockConn.EXPECT().Write(gomock.Any()).AnyTimes()
	mockConn.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (n int, err error) {
		data := q.Next()
		if data == nil {
			return 0, io.EOF
		}
		copy(p, data)
		return len(data), nil
	}).AnyTimes()

	switch code {
	case cmd.GetCode:
		mockHandler.EXPECT().HandleGet(gomock.Any(), ctx).AnyTimes().Return(nil)
	case cmd.SetCode:
		mockHandler.EXPECT().HandleSet(gomock.Any(), ctx).AnyTimes().Return(nil)
	case cmd.PingCode:
		mockHandler.EXPECT().HandlePing(gomock.Any()).AnyTimes().Return(nil)
	case cmd.SubCode:
		mockHandler.EXPECT().HandleSub(gomock.Any(), ctx).AnyTimes().Return(nil)
	case cmd.PubCode:
		mockHandler.EXPECT().HandlePub(gomock.Any(), ctx).AnyTimes().Return(nil)
	}
	processClient(mockConn, ctx, mockHandler, nil, nil)
}

func Test_processClientHandle(t *testing.T) {

	buildQueue := func(code string, values []string) *bytesQueue {
		q := new(bytesQueue)
		var str payload.String

		w := &bytes.Buffer{}
		str = payload.String(code)
		str.WriteTo(w)
		q.Add(w.Bytes())

		for _, v := range values {
			w = &bytes.Buffer{}
			str = payload.String(v)
			str.WriteTo(w)
			q.Add(w.Bytes())
		}
		return q
	}

	type args struct {
		ctx   *scontext.Context
		queue *bytesQueue
		code  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handle get success",
			args: args{
				ctx: &scontext.Context{
					Storage: func() storage.Storage {
						s := storage.NewDefaultStorage()
						s.Set("key", *payload.NewGeneral(payload.StringType, []byte("value")))
						return s
					}(),
					PubSub:  pubsub.NewPubSub(),
					Context: context.Background(),
				},
				code: cmd.GetCode,
				queue: func() *bytesQueue {
					return buildQueue(cmd.GetCode, []string{"key"})
				}(),
			},
		},
		{
			name: "handle set success",
			args: args{
				ctx: &scontext.Context{
					Storage: func() storage.Storage {
						s := storage.NewDefaultStorage()
						s.Set("key", *payload.NewGeneral(payload.StringType, []byte("value")))
						return s
					}(),
					PubSub:  pubsub.NewPubSub(),
					Context: context.Background(),
				},
				code: cmd.SetCode,
				queue: func() *bytesQueue {
					return buildQueue(cmd.SetCode, []string{"key", "value"})
				}(),
			},
		},
		{
			name: "handle pub success",
			args: args{
				ctx: &scontext.Context{
					Storage: storage.NewDefaultStorage(),
					PubSub:  pubsub.NewPubSub(),
					Context: context.Background(),
				},
				code: cmd.PubCode,
				queue: func() *bytesQueue {
					return buildQueue(cmd.PubCode, []string{"topic", "message"})
				}(),
			},
		},
		{
			name: "handle sub success",
			args: args{
				ctx: &scontext.Context{
					Storage: storage.NewDefaultStorage(),
					PubSub:  pubsub.NewPubSub(),
					Context: context.Background(),
				},
				code: cmd.SubCode,
				queue: func() *bytesQueue {
					return buildQueue(cmd.SubCode, []string{"topic"})
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandlerHelper(t, tt.args.ctx, tt.args.code, tt.args.queue)
		})
	}
}

func Test_getCommandFromConn(t *testing.T) {
	t.Skip()
	type args struct {
		cmdstr string
		conn   net.Conn
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "get ping command from conn",
			args: args{
				cmdstr: "get foo",
				conn:   mock.NewMockConn(gomock.NewController(t)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPayloadCommandFromConn(tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCommandFromConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCommandFromConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_restoreServer(t *testing.T) {
	type fields struct {
		config *Config
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "restore server success",
			fields: fields{
				config: &Config{
					PersistentPath: func() string {
						fileName := "./testdata/test-restore.save"
						_, err := os.Stat(fileName)
						if !os.IsNotExist(err) {
							return fileName
						}

						fileName = "../testdata/test-restore.save"
						_, err = os.Stat(fileName)
						if !os.IsNotExist(err) {
							return fileName
						}
						_, _ = os.CreateTemp(".", "")
						return ""
					}(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := handler.NewMockHandler(ctrl)
			handler.EXPECT().HandleSet(gomock.Any(), gomock.Any()).Times(4).Return(nil)
			s := &Server{
				handler: handler,
				config:  tt.fields.config,
			}
			s.restoreServer()
		})
	}
}
