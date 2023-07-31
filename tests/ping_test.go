package tests

import (
	"testing"

	common "github.com/panupakm/mini-redis/tests/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {

	_, c, close := common.SetUpServerClient(t)
	defer close()

	_, err := c.Ping("Yahoo")
	require.NoError(t, err)

	pong, err := c.ReadString()
	require.NoError(t, err)
	assert.Equal(t, "Yahoo", pong)
}
