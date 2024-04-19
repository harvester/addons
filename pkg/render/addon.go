package render

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

const (
	relativeTemplatePath = "../templates"
	defaultFileName      = "rancherd-22-addons.yaml"
)

type AddonResources struct {
	Resources []map[string]interface{} `json:"resources,omitempty"`
}

func Addon(templateSource, destPath string) error {
	contents, err := os.ReadFile(filepath.Join(templateSource, defaultFileName))
	if err != nil {
		return fmt.Errorf("error reading template file %s: %v", defaultFileName, err)
	}

	renderedContent, err := renderTemplate(contents)
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
		err = os.WriteFile(filepath.Join(destPath, fileName), addonContents, 0755)
		if err != nil {
			return fmt.Errorf("error writing addon file %s: %v", fileName, err)
		}
	}
	return nil
}

func renderTemplate(content []byte) ([]byte, error) {
	result := bytes.NewBufferString("")
	tmpl, err := template.New("").Parse(string(content))
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(result, nil)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
