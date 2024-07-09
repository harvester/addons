package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Addon(t *testing.T) {
	assert := require.New(t)
	tmpDir := os.TempDir()
	tmpPath, err := os.MkdirTemp(tmpDir, "rendered")
	assert.NoError(err)
	defer os.RemoveAll(tmpPath)
	err = Addon(relativeTemplatePath, tmpPath, filepath.Join(relativeVersionFilePath, defaultVersionFile))
	assert.NoError(err)
}

func Test_generate_version_info_map(t *testing.T) {
	assert := require.New(t)
	result, err := generate_version_info_map(filepath.Join(relativeVersionFilePath, defaultVersionFile))
	assert.NoError(err, "expected no error while reading version info map")
	_, ok := result["VM_IMPORT_CONTROLLER_CHART_VERSION"]
	assert.True(ok, "expected to find key VM_IMPORT_CONTROLLER_CHART_VERSION")
}

func Test_Template(t *testing.T) {
	assert := require.New(t)
	tmpDir := os.TempDir()
	tmpPath, err := os.MkdirTemp(tmpDir, "template")
	assert.NoError(err, "expected no error during tmp dir creation")
	defer os.RemoveAll(tmpPath)
	err = Template(relativeTemplatePath, tmpPath, filepath.Join(relativeVersionFilePath, defaultVersionFile))
	assert.NoError(err)
}
