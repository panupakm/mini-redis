package client

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/panupakm/miniredis/mock"
	"github.com/panupakm/miniredis/payload"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "new client",
			want: &Client{
				dial: &ImplDialer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestClient_ErrorConnect(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "no server with invalid address valid port",
			args: args{
				addr: "sadfs:23333",
			},
		},
		{
			name: "no server invalid address valid port",
			args: args{
				addr: "127.0.0.1:AFDF",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			err := c.Connect(tt.args.addr)
			assert.Error(t, err)
		})
	}
}

func TestClient_SuccessConnect(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "mock sersver with valid address valid port",
			args: args{
				addr: "127.0.0.1:9191",
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDialer := NewMockDialer(ctrl)
	mockConn := mock.NewMockConn(ctrl)
	mockDialer.EXPECT().Dial("tcp", "127.0.0.1:9191").Return(mockConn, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				dial: mockDialer,
			}
			err := c.Connect(tt.args.addr)
			assert.NoError(t, err)
		})
	}
}

func TestClient_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConn := mock.NewMockConn(ctrl)

	type fields struct {
		conn net.Conn
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "close valid connection",
			fields: fields{
				conn: mockConn,
				addr: "127.0.0.1:1999",
			},
			wantErr: false,
		},
		{
			name: "close null connection",
			fields: fields{
				conn: nil,
				addr: "127.0.0.1:1999",
			},
			wantErr: true,
		},
	}

	mockConn.EXPECT().Close().Times(1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				conn: tt.fields.conn,
				addr: tt.fields.addr,
			}
			err := c.Close()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func configureMockConnResult(code uint16, msg string, ctrl *gomock.Controller) *mock.MockConn {
	mockResultBytes := func() []byte {
		rsPayload := payload.Result{
			Code:   code,
			Length: uint32(len(msg)),
			Typ:    payload.StringType,
			Buffer: []byte(msg),
		}.Bytes()
		bytes := []byte{byte(payload.ResultType)}
		lenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBytes, uint32(len(rsPayload)))
		bytes = append(bytes, lenBytes...)
		return append(bytes, rsPayload...)
	}()

	mockConn := mock.NewMockConn(ctrl)

	dd := make(map[net.Conn]string)
	dd[mockConn] = string(mockResultBytes)

	mockConn.EXPECT().Write(gomock.Any()).AnyTimes()
	mockConn.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		readByteCount := copy(b, mockResultBytes)
		mockResultBytes = mockResultBytes[readByteCount:]
		return readByteCount, nil
	}).AnyTimes()

	return mockConn
}

func TestClient_Ping(t *testing.T) {
	type fields struct {
		addr string
	}
	type args struct {
		msg  string
		code uint16
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      chan ResultChannel
		callTimes int
		wantErr   bool
	}{
		{
			name: "ping valid connection",
			fields: fields{
				addr: "127.0.0.1:1999",
			},
			args: args{
				msg:  "pong",
				code: 0,
			},
			callTimes: 6,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := configureMockConnResult(tt.args.code, tt.args.msg, ctrl)
			c := &Client{
				conn: mockConn,
				addr: tt.fields.addr,
			}

			got, err := c.Ping(tt.args.msg)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			rs := <-got
			assert.NoError(t, rs.Err)
			assert.Equal(t, tt.args.msg, rs.DataAsString())
		})
	}
}

func TestClient_SetString(t *testing.T) {
	type fields struct {
		addr string
	}
	type args struct {
		key   string
		value string
		code  uint16
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    chan ResultChannel
		wantMsg string
		wantErr bool
	}{

		{
			name: "set string valid connection",
			fields: fields{
				addr: "127.0.0.1:1999",
			},
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: false,
			wantMsg: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := configureMockConnResult(tt.args.code, tt.wantMsg, ctrl)
			c := &Client{
				conn: mockConn,
				addr: tt.fields.addr,
			}
			got, err := c.SetString(tt.args.key, tt.args.value)
			assert.NoError(t, err)
			assert.NotNil(t, got)

			rs := <-got

			assert.NoError(t, rs.Err)
			assert.Equal(t, tt.wantMsg, rs.DataAsString())
		})
	}
}

func TestClient_Get(t *testing.T) {
	type fields struct {
		addr string
	}
	type args struct {
		key  string
		code uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    chan ResultChannel
		wantMsg string
		wantErr bool
	}{
		{
			name: "get with valid key",
			fields: fields{
				addr: "127.0.0.1:1999",
			},
			args: args{
				key: "key",
			},
			wantErr: false,
			wantMsg: "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := configureMockConnResult(tt.args.code, tt.wantMsg, ctrl)
			c := &Client{
				conn: mockConn,
				addr: tt.fields.addr,
			}

			got, err := c.Get(tt.args.key)
			assert.NoError(t, err)
			assert.NotNil(t, got)

			rs := <-got

			assert.NoError(t, rs.Err)
			assert.Equal(t, tt.wantMsg, rs.DataAsString())
		})
	}
}

func TestClient_Sub(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
		dial Dialer
	}
	type args struct {
		topic string
		code  uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Subsriber
		wantMsg string
	}{
		{
			name: "sub valid topic",
			fields: fields{
				addr: "127.0.0.1:1999",
			},
			args: args{
				topic: "topic",
				code:  0,
			},
			wantMsg: "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := configureMockConnResult(tt.args.code, tt.wantMsg, ctrl)
			c := &Client{
				conn: mockConn,
				addr: tt.fields.addr,
			}
			sub, ch, err := c.Sub(tt.args.topic)
			assert.NoError(t, err)
			assert.NotNil(t, sub)
			rs := <-ch
			assert.NoError(t, rs.Err)
			assert.Equal(t, tt.wantMsg, rs.DataAsString())
		})
	}
}

func TestClient_PubString(t *testing.T) {
	type fields struct {
		conn net.Conn
		addr string
	}
	type args struct {
		topic string
		msg   string
		code  uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMsg string
	}{
		{
			name: "publish string message",
			fields: fields{
				addr: "127.0.0.1:1999",
			},
			args: args{
				topic: "topic",
				msg:   "msg",
			},
			wantMsg: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := configureMockConnResult(tt.args.code, tt.wantMsg, ctrl)
			c := &Client{
				conn: mockConn,
				addr: tt.fields.addr,
			}
			got, err := c.PubString(tt.args.topic, tt.args.msg)
			assert.NoError(t, err)
			rs := <-got
			assert.NoError(t, rs.Err)
			assert.Equal(t, tt.wantMsg, rs.DataAsString())
		})
	}
}
