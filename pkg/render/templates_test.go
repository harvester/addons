package render

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

const sourcePath = "../templates/"

func Test_Templates(t *testing.T) {
	type FileComparison struct {
		Actual   []byte
		Expected []byte
	}

	testData := make(map[string]FileComparison)

	assert := require.New(t)
	tmpDir := os.TempDir()
	tmpPath, err := os.MkdirTemp(tmpDir, "addons")
	assert.NoError(err)
	defer os.RemoveAll(tmpPath)
	err = Template(tmpPath)
	assert.NoError(err)

	// walk target dir and read templates
	err = filepath.Walk(tmpPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", info.Name(), err)
			}
			testData[info.Name()] = FileComparison{
				Actual: content,
			}
		}
		return nil
	})

	// walk source dir and read source templates
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", info.Name(), err)
			}
			val, ok := testData[info.Name()]
			if !ok {
				return fmt.Errorf("did not find file %s populated in test data", info.Name())
			}
			val.Expected = content
			testData[info.Name()] = val
		}
		return nil
	})
	assert.NoError(err)

	// compare files are equal
	for name, fileContent := range testData {
		t.Logf("comparing file %s", name)
		assert.True(bytes.Equal(fileContent.Actual, fileContent.Expected), fmt.Sprintf("expected and actual contents do not match for file %s", name))
	}
}
