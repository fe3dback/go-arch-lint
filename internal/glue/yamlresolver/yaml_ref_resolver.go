package yamlresolver

import (
	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type YamlReferenceResolver struct {
	sourceCode     []byte
	sourceFilePath string
	parser         AnnotationParser
}

func NewYamlReferenceResolver(
	sourceCode []byte,
	sourceFilePath string,
	parser AnnotationParser,
) *YamlReferenceResolver {
	return &YamlReferenceResolver{
		sourceCode:     sourceCode,
		sourceFilePath: sourceFilePath,
		parser:         parser,
	}
}

func (r *YamlReferenceResolver) Resolve(yamlPath string) (ref speca.Reference) {
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

	annotation, err := path.AnnotateSource(r.sourceCode, false)
	if err != nil {
		return speca.NewEmptyReference()
	}

	line, pos, err := r.parser.Parse(string(annotation))
	if err != nil {
		return speca.NewEmptyReference()
	}

	return speca.NewReference(
		r.sourceFilePath,
		line,
		pos,
	)
}
