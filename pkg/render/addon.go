package render

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

const (
	relativeTemplatePath    = "../templates"
	defaultFileName         = "rancherd-22-addons.yaml"
	relativeVersionFilePath = "../../"
	defaultVersionFile      = "version_info"
	imageSuffix             = "IMAGE"
)

type AddonResources struct {
	Resources []map[string]interface{} `json:"resources,omitempty"`
}

func Addon(templateSource, destPath, versionFilePath string) error {

	tmpDir := os.TempDir()
	tmpPath, err := os.MkdirTemp(tmpDir, "rendered")
	if err != nil {
		return fmt.Errorf("error creating temp addon-template file: %v", err)
	}
	defer os.RemoveAll(tmpPath)

	err = Template(templateSource, tmpPath, versionFilePath)
	if err != nil {
		return fmt.Errorf("error generating temp file: %v", err)
	}

	// read temporary template file to generate rendered addons
	contents, err := os.ReadFile(filepath.Join(tmpPath, defaultFileName))
	if err != nil {
		return fmt.Errorf("error reading template file %s: %v", defaultFileName, err)
	}

	tmpl, err := template.New("").Parse(string(contents))
	if err != nil {
		return err
	}

	renderedContent, err := renderTemplate(tmpl, versionFilePath)
	if err != nil {
		return fmt.Errorf("error rendering template: %v", err)
	}

	//split rendered template into individual files
	resources := &AddonResources{}
	err = yaml.Unmarshal(renderedContent, resources)
	if err != nil {
		return fmt.Errorf("error unmarshalling resources: %v", err)
	}

	for _, v := range resources.Resources {
		metadata, ok := v["metadata"]
		if !ok {
			logrus.Errorf("skipping resource since metadata is missing: %v", v)
			continue
		}
		metadataMap, ok := metadata.(map[string]interface{})
		if !ok {
			logrus.Errorf("skipping resource as unable to assert metadata to map[string]interface{}: %v", v)
		}

		name, ok := metadataMap["name"]
		if !ok {
			logrus.Errorf("skipping resource since name is missing in metadata: %v", v)
			continue
		}
		addonContents, err := yaml.Marshal(v)
		if err != nil {
			return fmt.Errorf("error marshalling map to addon contents for addon %s: %v", name, err)
		}
		fileName := fmt.Sprintf("%s.yaml", name)
		err = os.WriteFile(filepath.Join(destPath, fileName), addonContents, 0644)
		if err != nil {
			return fmt.Errorf("error writing addon file %s: %v", fileName, err)
		}
	}
	return nil
}

func renderTemplate(tmpl *template.Template, versionFilePath string) ([]byte, error) {
	envMap, err := generate_version_info_map(versionFilePath)
	if err != nil {
		return nil, err
	}

	result := bytes.NewBufferString("")
	err = tmpl.Execute(result, envMap)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

// generate_version_info_map reads the version_info file in project root
// and renders a raw template which can be used for subsequent processing
func generate_version_info_map(versionFilePath string) (map[string]string, error) {
	versionFile, err := os.Open(versionFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening version file: %v", err)
	}
	defer versionFile.Close()
	versionScanner := bufio.NewScanner(versionFile)

	var lines []string
	for versionScanner.Scan() {
		lines = append(lines, versionScanner.Text())
	}

	result := make(map[string]string, len(lines))
	for _, v := range lines {
		fields := strings.Split(v, "=")
		if len(fields) == 2 { // need check to ignore hash=bang directive in version info file
			if strings.Contains(fields[0], imageSuffix) {
				// in case of image names we only need tag info and not full image name
				images := strings.Split(fields[1], ":")
				result[fields[0]] = strings.Trim(images[1], "\"")

			} else {
				// in case of a helm chart there the string will not be of format image:tag
				// but will only contain chart version
				// for example NVIDIA_DRIVER_RUNTIME_CHART_VERSION="0.1.1"
				result[fields[0]] = fields[1]
			}
		}
	}
	return result, nil
}

func Template(templateSource, destPath, versionFilePath string) error {
	contents, err := os.ReadFile(filepath.Join(templateSource, defaultFileName))
	if err != nil {
		return fmt.Errorf("error reading template file %s: %v", defaultFileName, err)
	}

	tmpl, err := template.New("").Delims("<<", ">>").Parse(string(contents))
	if err != nil {
		return err
	}
	renderedContent, err := renderTemplate(tmpl, versionFilePath)
	if err != nil {
		return fmt.Errorf("error rendering template: %v", err)
	}

	err = os.WriteFile(filepath.Join(destPath, defaultFileName), renderedContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing addon file %s: %v", defaultFileName, err)
	}

	return nil
}
