package config

import (
	"fmt"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Reader struct {
	yamlReader FileReader
}

func NewReader(
	yamlReader FileReader,
) *Reader {
	return &Reader{
		yamlReader: yamlReader,
	}
}

func (r *Reader) Read(path models.PathAbsolute) (models.Config, error) {
	ext := filepath.Ext(string(path))

	switch ext {
	case ".yml", ".yaml":
		return r.yamlReader.ReadFile(path)
	default:
		return models.Config{}, fmt.Errorf("unknown config file '%s' ext: %s", path, ext)
	}
}
