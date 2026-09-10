package render

import (
	"os"
	"path/filepath"
	"strings"
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
	assert.Contains(result, "VM_IMPORT_CONTROLLER_CHART_VERSION")

	expectedImages := map[string]struct {
		repository string
		tag        string
	}{
		"VM_IMPORT_CONTROLLER":  {repository: "rancher/harvester-vm-import-controller", tag: "main-head"},
		"PCIDEVICES_CONTROLLER": {repository: "rancher/harvester-pcidevices", tag: "master-head"},
		"HARVESTER_SEEDER":      {repository: "rancher/harvester-seeder", tag: "main-head"},
		"NVIDIA_DRIVER_TOOLKIT": {repository: "rancher/harvester-nvidia-driver-toolkit", tag: "sle-micro-head"},
		"HARVESTER_EVENTROUTER": {repository: "rancher/harvester-eventrouter", tag: "master-head"},
		"KUBEOVN_OPERATOR":      {repository: "rancher/harvester-kubeovn-operator", tag: "v1.16.2-rc6"},
		"DESCHEDULER":           {repository: "registry.k8s.io/descheduler/descheduler", tag: "v0.36.0"},
	}
	for key, expected := range expectedImages {
		assert.Equal(expected.repository, result[key+"_IMAGE_REPOSITORY"])
		assert.Equal(expected.tag, result[key+"_IMAGE_TAG"])
		assert.NotContains(result, key+"_IMAGE")
	}
}

func Test_generate_version_info_mapPreservesRegistryPort(t *testing.T) {
	assert := require.New(t)
	versionFile := filepath.Join(t.TempDir(), "version_info")
	err := os.WriteFile(versionFile, []byte("TEST_IMAGE_REPOSITORY=\"registry.example.com:5000/team/image\"\nTEST_IMAGE_TAG=\"v1.2.3\"\n"), 0600)
	assert.NoError(err)

	result, err := generate_version_info_map(versionFile)
	assert.NoError(err)
	assert.Equal("registry.example.com:5000/team/image", result["TEST_IMAGE_REPOSITORY"])
	assert.Equal("v1.2.3", result["TEST_IMAGE_TAG"])
}

func Test_Template(t *testing.T) {
	assert := require.New(t)
	tmpDir := os.TempDir()
	tmpPath, err := os.MkdirTemp(tmpDir, "template")
	assert.NoError(err, "expected no error during tmp dir creation")
	defer os.RemoveAll(tmpPath)
	err = Template(relativeTemplatePath, tmpPath, filepath.Join(relativeVersionFilePath, defaultVersionFile))
	assert.NoError(err)

	content, err := os.ReadFile(filepath.Join(tmpPath, defaultFileName))
	assert.NoError(err)
	rendered := string(content)
	for _, expected := range []string{
		"repository: rancher/harvester-vm-import-controller",
		"repository: rancher/harvester-pcidevices",
		"repository: rancher/harvester-seeder",
		"repository: rancher/harvester-nvidia-driver-toolkit",
		"repository: rancher/harvester-kubeovn-operator",
		"repository: registry.k8s.io/descheduler/descheduler",
	} {
		assert.True(strings.Contains(rendered, expected), "expected rendered template to contain %q", expected)
	}
}
