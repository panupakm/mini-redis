package tests

import (
	"testing"

	common "github.com/panupakm/miniredis/tests/internal"
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

			_, await, err := c.Sub(tt.args.topic)
			result := <-await
			require.NoError(t, result.Err)
			require.Equal(t, tt.args.result, string(result.Buffer))

			require.NoError(t, err)
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
			csub1, csubclose1 := common.SetUpClient(t, port)
			defer csubclose1()
			csub2, csubclose2 := common.SetUpClient(t, port)
			defer csubclose2()

			cpub, cpubclose := common.SetUpClient(t, port)
			defer cpubclose()

			//two clients subscribing to the same topic
			subscriber1, await, err := csub1.Sub(tt.args.topic)
			require.NoError(t, err)
			result := <-await
			require.NoError(t, result.Err)
			require.Equal(t, tt.args.result, string(result.Buffer))

			subscriber2, await, err := csub2.Sub(tt.args.topic)
			require.NoError(t, err)
			result = <-await
			require.NoError(t, result.Err)
			require.Equal(t, tt.args.result, string(result.Buffer))

			//publish string message
			_, err = cpub.PubString(tt.args.topic, tt.args.msg)
			require.NoError(t, err)

			//wait for message
			msg, _ := subscriber1.NextMessage()
			str, _ := msg.AsString()
			require.Equal(t, tt.args.msg, str)

			msg, _ = subscriber2.NextMessage()
			str, _ = msg.AsString()
			require.Equal(t, tt.args.msg, str)
		})
	}
}
