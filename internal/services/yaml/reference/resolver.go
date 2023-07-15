package reference

import (
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

type Resolver struct {
}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Resolve(filePath string, yamlPath string) (ref models.Reference) {
	defer func() {
		if data := recover(); data != nil {
			ref = speca.NewEmptyReference()
			ref.Hint = fmt.Sprintf("%s", data)
			return
		}
	}()

	sourceCode, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to provide source code of archfile: %v", err))
	}

	path, err := yaml.PathString(yamlPath)
	if err != nil {
		return speca.NewEmptyReference()
	}

	file, err := parser.ParseBytes(sourceCode, 0)
	if err != nil {
		return speca.NewEmptyReference()
	}

	node, err := path.FilterFile(file)
	if err != nil {
		return speca.NewEmptyReference()
	}

	pos := node.GetToken().Position

	return speca.NewReference(
		filePath,
		pos.Line,
		pos.Column,
	)
}
