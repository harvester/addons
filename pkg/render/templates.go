package render

import (
	"github.com/harvester/addons/pkg/data"
	"log"
	"os"
	"path/filepath"
)

func Template(path string) error {
	for _, v := range data.AssetNames() {
		content, err := data.Asset(v)
		if err != nil {
			return err
		}
		_, fileName := filepath.Split(v)
		log.Printf("writing template file %s", fileName)
		if err := os.WriteFile(filepath.Join(path, fileName), content, 0755); err != nil {
			return err
		}
	}
	return nil
}
