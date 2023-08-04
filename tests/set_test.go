package tests

import (
	"testing"

	common "github.com/panupakm/miniredis/tests/internal"
	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	type args struct {
		key   string
		value string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set value as string",
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: false,
		},
	}
	port := "9989"
	_ = common.SetUpServer(t, port)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, close := common.SetUpClient(t, port)
			defer close()

			ch, err := c.SetString(tt.args.key, tt.args.value)
			require.NoError(t, err)
			r := <-ch
			require.NoError(t, r.Err)
			require.Equal(t, "OK", string(r.Buffer))

			ch, err = c.Get(tt.args.key)
			require.NoError(t, err)
			r = <-ch
			require.NoError(t, r.Err)
			require.Equal(t, tt.args.value, string(r.Buffer))
		})
	}
}
