package tests

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	common "github.com/panupakm/miniredis/tests/internal"
	"github.com/stretchr/testify/assert"
)

func TestConcurrency(t *testing.T) {
	type args struct {
		msg         string
		topic       string
		clientCount int
		loopCount   int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "stress test muliple clients",
			args: args{
				msg:         "PING",
				topic:       "test",
				clientCount: 10,
				loopCount:   10,
			},
			wantErr: false,
		},
	}
	port := uint(8100)
	_ = common.SetUpServer(t, port, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < tt.args.clientCount; i++ {
				wg.Add(1)
				go func(no int) {
					client := common.SetUpClient(t, 8100)
					defer client.Close()

					for j := 0; j < tt.args.loopCount; j++ {
						rsch, _ := client.Get(tt.args.topic)
						rs := <-rsch
						assert.NoError(t, rs.Err)
						nextUint, _ := strconv.ParseUint(string(rs.Result.Buffer), 10, 0)
						rsch, _ = client.SetString(tt.args.topic, fmt.Sprint(nextUint+1))
						rs = <-rsch
						assert.NoError(t, rs.Err)
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
			client := common.SetUpClient(t, 8100)
			defer client.Close()
			rsch, _ := client.Get(tt.args.topic)
			rs := <-rsch
			assert.NoError(t, rs.Err)
		})
	}
}
