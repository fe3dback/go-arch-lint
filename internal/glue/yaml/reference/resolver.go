package reference

import (
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Resolver struct {
	sourceCode     []byte
	sourceFilePath string
}

func NewResolver(
	sourceCode []byte,
	sourceFilePath string,
) *Resolver {
	return &Resolver{
		sourceCode:     sourceCode,
		sourceFilePath: sourceFilePath,
	}
}

func (r *Resolver) Resolve(yamlPath string) (ref models.Reference) {
	defer func() {
		if data := recover(); data != nil {
			ref = speca.NewEmptyReference()
			return
		}
	}()

	path, err := yaml.PathString(yamlPath)
	if err != nil {
		return speca.NewEmptyReference()
	}

	file, err := parser.ParseBytes(r.sourceCode, 0)
	if err != nil {
		return speca.NewEmptyReference()
	}

	node, err := path.FilterFile(file)
	if err != nil {
		return speca.NewEmptyReference()
	}

	pos := node.GetToken().Position

	return speca.NewReference(
		r.sourceFilePath,
		pos.Line,
		pos.Column,
	)
}
