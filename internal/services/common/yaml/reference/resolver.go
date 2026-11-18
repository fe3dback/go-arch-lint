package reference

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fe3dback/go-yaml"
	"github.com/fe3dback/go-yaml/parser"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
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

	// archfile written on windows usually has CRLF line endings, yaml parser
	// count '\r' as part of line and return shifted line/column,
	// so all notices will point to wrong place of archfile
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

	r.cache[filePath] = content
	return content
}
