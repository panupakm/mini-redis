package request

import (
	"bytes"
	"io"
	"testing"

	"github.com/panupakm/miniredis/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makePubPayloadReader[T payload.Payload](topic string, p T) io.Reader {
	buf := new(bytes.Buffer)
	s := payload.String(topic)
	s.WriteTo(buf)
	p.WriteTo(buf)
	return bytes.NewReader(buf.Bytes())
}

func makePubPayloadStringWriter(topic string, msg string) []byte {
	buf := new(bytes.Buffer)
	s := payload.String(PubCode)
	s.WriteTo(buf)
	s = payload.String(topic)
	s.WriteTo(buf)
	s = payload.String(msg)
	s.WriteTo(buf)
	return buf.Bytes()
}

func TestPubReadFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	const mintedTopic = "minted"
	const mintedMsg = "{tokenId: 123456789012345678}"

	tests := []struct {
		name string
		args args
		want *Pub
	}{
		{
			name: "string topic and string msg",
			args: args{
				r: func() io.Reader {
					p := payload.String(mintedMsg)
					return makePubPayloadReader(mintedTopic, &p)
				}(),
			},
			want: &Pub{
				Topic: mintedTopic,
				Typ:   payload.StringType,
				Data:  []byte(mintedMsg),
				Len:   uint32(len(mintedMsg)),
			},
		},
		{
			name: "string topic and empty string msg",
			args: args{
				r: func() io.Reader {
					p := payload.String("")
					return makePubPayloadReader(mintedTopic, &p)
				}(),
			},
			want: &Pub{
				Topic: mintedTopic,
				Typ:   payload.StringType,
				Data:  []byte(""),
				Len:   uint32(len("")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PubReadFrom(tt.args.r)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPub_String(t *testing.T) {
	type fields struct {
		Topic string
		Typ   payload.ValueType
		Len   uint32
		Data  []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string topic and string msg",
			fields: fields{
				Topic: "greeting",
				Typ:   payload.StringType,
				Data:  []byte("msg"),
				Len:   uint32(len("msg")),
			},
			want: "pub topic:greeting",
		},
		{
			name: "string topic and empty string msg",
			fields: fields{
				Topic: "greeting",
				Typ:   payload.StringType,
				Data:  []byte(""),
				Len:   uint32(len("")),
			},
			want: "pub topic:greeting",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Pub{
				Topic: tt.fields.Topic,
				Typ:   tt.fields.Typ,
				Len:   tt.fields.Len,
				Data:  tt.fields.Data,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("Pub.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubStringTo(t *testing.T) {
	type args struct {
		topic string
		msg   string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "string topic and string msg",
			args: args{
				topic: "greeting",
				msg:   "sawasdee",
			},
			want: makePubPayloadStringWriter("greeting", "sawasdee"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := PubStringTo(w, tt.args.topic, tt.args.msg)
			require.NoError(t, err)
			assert.Equal(t, tt.want, w.Bytes())
		})
	}
}

func TestPub_Bytes(t *testing.T) {
	type fields struct {
		Topic string
		Typ   payload.ValueType
		Len   uint32
		Data  []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "valid topic valid string msg",
			fields: fields{
				Topic: "greeting",
				Typ:   payload.StringType,
				Data:  []byte("hello world"),
				Len:   uint32(len("hello world")),
			},
			want: makePubPayloadStringWriter("greeting", "hello world"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Pub{
				Topic: tt.fields.Topic,
				Typ:   tt.fields.Typ,
				Len:   tt.fields.Len,
				Data:  tt.fields.Data,
			}
			got := s.Bytes()
			assert.Equal(t, tt.want, got)
		})
	}
}
