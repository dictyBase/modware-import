package manifest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomAppName(t *testing.T) {
	assert := require.New(t)
	name, err := RandomAppName(10)
	assert.NoError(err, "expect no error from generating app name")
	assert.Lenf(name, 10, "expect app name of length %d got %d", 10, len(name))
}
