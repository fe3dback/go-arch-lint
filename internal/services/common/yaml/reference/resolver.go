package reference

import (
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

type Resolver struct {
	cache map[string][]byte
}

func NewResolver() *Resolver {
	return &Resolver{
		cache: map[string][]byte{},
	}
}

func (r *Resolver) Resolve(filePath string, yamlPath string) (ref common.Reference) {
	defer func() {
		if data := recover(); data != nil {
			ref = common.NewEmptyReference()
			return
		}
	}()

	sourceCode := r.fileSource(filePath)

	path, err := yaml.PathString(yamlPath)
	if err != nil {
		return common.NewEmptyReference()
	}

	file, err := parser.ParseBytes(sourceCode, 0)
	if err != nil {
		return common.NewEmptyReference()
	}

	node, err := path.FilterFile(file)
	if err != nil {
		return common.NewEmptyReference()
	}

	pos := node.GetToken().Position

	return common.NewReferenceSingleLine(
		filePath,
		pos.Line,
		pos.Column,
	)
}

func (r *Resolver) fileSource(filePath string) []byte {
	if content, exist := r.cache[filePath]; exist {
		return content
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to provide source code of archfile: %v", err))
	}

	r.cache[filePath] = content
	return content
}
