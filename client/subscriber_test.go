package client

import (
	"encoding/binary"
	"testing"

	"github.com/panupakm/miniredis/mock"
	"github.com/panupakm/miniredis/payload"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func configureMockConnForSub(msg string, ctrl *gomock.Controller) *mock.MockConn {
	mockSubMsgBytes := func() []byte {
		bytes := []byte{byte(payload.SubMsgType)}

		bytes = append(bytes, byte(payload.StringType))
		lenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBytes, uint32(len(msg)))
		bytes = append(bytes, lenBytes...)
		return append(bytes, []byte(msg)...)
	}()

	mockConn := mock.NewMockConn(ctrl)
	mockConn.EXPECT().Write(gomock.Any()).AnyTimes()
	mockConn.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		readByteCount := copy(b, mockSubMsgBytes)
		mockSubMsgBytes = mockSubMsgBytes[readByteCount:]
		return readByteCount, nil
	}).AnyTimes()

	return mockConn
}

func TestSubsriber_NextMessage(t *testing.T) {
	type fields struct {
		messages chan payload.SubMsg
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:   "Read next message",
			fields: fields{},
			want:   "message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := configureMockConnForSub(tt.want, ctrl)
			s := &Subsriber{
				messages: tt.fields.messages,
				conn:     mockConn,
			}
			got, err := s.NextMessage()
			assert.NoError(t, err)
			msg, _ := got.AsString()
			assert.Equal(t, tt.want, msg)
		})
	}
}
