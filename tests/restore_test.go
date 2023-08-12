package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/panupakm/miniredis/server"
	common "github.com/panupakm/miniredis/tests/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerRestore(t *testing.T) {

	type pair struct {
		key string
		val string
	}
	type args struct {
		set []pair
	}
	tests := []struct {
		name string
		args args
		want []pair
	}{
		{
			name: "restoring server",
			args: args{
				set: []pair{
					{
						key: "key1",
						val: "value1",
					},
					{
						key: "key2",
						val: "value2",
					},
					{
						key: "key3",
						val: "value3",
					},
					{
						key: "key1",
						val: "value11",
					},
				},
			},
			want: []pair{
				{
					key: "key2",
					val: "value2",
				},
				{
					key: "key3",
					val: "value3",
				},
				{
					key: "key1",
					val: "value11",
				},
			},
		},
	}

	tempDir, err := os.MkdirTemp(".", "restore-test")
	require.NoError(t, err)
	// defer os.RemoveAll(tempDir)

	const startPort = 9801
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savePath := filepath.Join(tempDir, "test-persistent-1.save")
			s1 := common.SetUpServer(t, startPort, &server.Config{
				PersistentPath: savePath,
			})
			defer func() {
				if s1 != nil {
					s1.Close()
				}
			}()
			c := common.SetUpClient(t, startPort)
			defer c.Close()

			for _, p := range tt.args.set {
				ch, err := c.SetString(p.key, p.val)
				assert.NoError(t, err)
				rs := <-ch
				assert.Equal(t, "OK", string(rs.Buffer))
			}
			s1.Close()
			s1 = nil

			// Start a new server to restore data from shutdown server
			s2 := common.SetUpServer(t, startPort+1, &server.Config{
				PersistentPath: savePath,
			})
			defer s2.Close()
			time.Sleep(1 * time.Second)

			for _, p := range tt.want {
				ch, err := c.Get(p.key)
				assert.NoError(t, err)
				rs := <-ch
				assert.Equal(t, p.val, string(rs.Buffer))
			}
		})
	}
}
