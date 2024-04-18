package render

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	relativeTemplatePath = "../templates"
)

func Addon(templateSource, destPath string) error {
	return filepath.Walk(templateSource, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			templateContent, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", info.Name(), err)
			}
			generatedContent, err := renderTemplate(templateContent)
			if err != nil {
				return fmt.Errorf("error rendering template %s: %v", info.Name(), err)
			}
			err = os.WriteFile(filepath.Join(destPath, info.Name()), generatedContent, 0755)
			if err != nil {
				return fmt.Errorf("error writing addon %s: %v", info.Name(), err)
			}
			log.Printf("generated addon file %s", info.Name())
		}
		return nil
	})
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
