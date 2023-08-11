package tests

import (
	"testing"

	common "github.com/panupakm/miniredis/tests/internal"
	"github.com/stretchr/testify/require"
)

func TestServerRestore(t *testing.T) {

	type pair struct {
		key string
		val string
	}
	type args struct {
		pubsub []pair
		port   uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set value as string",
			args: args{
				pubsub: []pair{
					{
						key: "topic1",
						val: "value1",
					},
				},
			},
			wantErr: false,
		},
	}

	startPort := 9800
	for i := range tests {
		tests[i].args.port = uint(startPort + i)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = common.SetUpServer(t, tt.args.port)
			c, close := common.SetUpClient(t, tt.args.port)
			defer close()

			for _, p := range tt.args.pubsub {
				_, await, err := c.Sub(p.key)
				require.NoError(t, err)
				result := <-await
				require.NoError(t, result.Err)
			}
		})
	}
}
