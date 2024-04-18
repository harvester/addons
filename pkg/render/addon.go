package render

import (
	"bytes"
	"github.com/harvester/addons/pkg/data"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func Addon(path string) error {
	for _, v := range data.AssetNames() {
		content, err := data.Asset(v)
		if err != nil {
			return err
		}
		_, fileName := filepath.Split(v)

		renderedContent, err := renderTemplate(content)
		if err != nil {
			return err
		}
		log.Printf("writing rendered addon file %s", fileName)
		if err := os.WriteFile(filepath.Join(path, fileName), renderedContent, 0755); err != nil {
			return err
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
