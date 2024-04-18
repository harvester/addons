package render

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_Addon(t *testing.T) {
	assert := require.New(t)
	tmpDir := os.TempDir()
	tmpPath, err := os.MkdirTemp(tmpDir, "rendered")
	assert.NoError(err)
	defer os.RemoveAll(tmpPath)
	err = Addon(tmpPath)
	assert.NoError(err)
}
