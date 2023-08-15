package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGeneral(t *testing.T) {
	type args struct {
		typ  ValueType
		buff []byte
	}
	tests := []struct {
		name string
		args args
		want *General
	}{
		{
			name: "create general payload from string payload",
			args: args{
				typ:  StringType,
				buff: []byte("hello"),
			},
			want: NewGeneral(StringType, []byte("hello")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGeneral(tt.args.typ, tt.args.buff)
			assert.Equal(t, tt.want, got)
		})
	}
}
