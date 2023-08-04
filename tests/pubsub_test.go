package tests

import (
	"testing"

	common "github.com/panupakm/miniredis/tests/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSub(t *testing.T) {
	type args struct {
		topic  string
		result string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "subscribe new topic",
			args: args{
				topic:  "hello",
				result: "OK",
			},
			wantErr: false,
		},
		{
			name: "subscribe same topic",
			args: args{
				topic:  "hello",
				result: "OK",
			},
			wantErr: false,
		},
	}

	port := "9990"

	_ = common.SetUpServer(t, port)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, close := common.SetUpClient(t, port)
			defer close()

			ch, err := c.Sub(tt.args.topic)
			require.NoError(t, err)
			r := <-ch
			require.NoError(t, r.Err)
			assert.Equal(t, tt.args.result, r.Str)
		})
	}
}

func TestPubToExistingTopic(t *testing.T) {
	type args struct {
		topic  string
		msg    string
		result string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pub to existing topic",
			args: args{
				topic:  "greeting",
				msg:    "sawasdee",
				result: "OK",
			},
			wantErr: false,
		},
	}

	port := "9991"

	_ = common.SetUpServer(t, port)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csub, csubclose := common.SetUpClient(t, port)
			defer csubclose()

			cpub, cpubclose := common.SetUpClient(t, port)
			defer cpubclose()

			chsub, err := csub.Sub(tt.args.topic)
			require.NoError(t, err)
			rs := <-chsub
			require.NoError(t, rs.Err)
			require.Equal(t, tt.args.result, string(rs.Str))

			chpub, err := cpub.PubString(tt.args.topic, tt.args.msg)
			require.NoError(t, err)
			rs = <-chpub
			require.NoError(t, rs.Err)
			require.Equal(t, uint16(0), rs.Code)
			require.Equal(t, "OK", string(rs.Buffer))
		})
	}
}
