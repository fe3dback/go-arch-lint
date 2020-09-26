package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/spec/annotate"
	"github.com/goccy/go-yaml"
)

type (
	YamlAnnotatedWarning struct {
		Text       string
		Path       string
		Line       int
		Offset     int
		SourceCode *YamlAnnotatedWarningSource `json:"-"`
	}

	YamlAnnotatedWarningSource struct {
		FormatText          []byte
		FormatTextHighlight []byte
	}

	AnnotatedValidator struct {
		validator  *ArchFileValidator
		sourceCode []byte
	}
)

func NewAnnotatedValidator(validator *ArchFileValidator, sourceCode []byte) *AnnotatedValidator {
	return &AnnotatedValidator{
		validator:  validator,
		sourceCode: sourceCode,
	}
}

func (av *AnnotatedValidator) Validate() ([]YamlAnnotatedWarning, error) {
	return av.annotateWarnings(
		av.validator.Validate(),
	)
}

func (av *AnnotatedValidator) annotateWarnings(warnings []Warning) ([]YamlAnnotatedWarning, error) {
	annotatedWarnings := make([]YamlAnnotatedWarning, 0)

	for _, warning := range warnings {
		path := warning.Path()

		sourceLine, err := yaml.PathString(path)
		if err != nil {
			return nil, fmt.Errorf("failed check '%s': %v", path, err)
		}

		textSource, err := sourceLine.AnnotateSource(av.sourceCode, false)
		if err != nil {
			return nil, fmt.Errorf("failed annotate '%s': %v", path, err)
		}

		highlightSource, err := sourceLine.AnnotateSource(av.sourceCode, true)
		if err != nil {
			return nil, fmt.Errorf("failed annotate '%s': %v", path, err)
		}

		sourceMarker := annotate.ParseSourceError(string(textSource))

		annotatedWarnings = append(annotatedWarnings, YamlAnnotatedWarning{
			Text:   warning.Warning().Error(),
			Path:   path,
			Line:   sourceMarker.Line,
			Offset: sourceMarker.Pos,
			SourceCode: &YamlAnnotatedWarningSource{
				FormatText:          textSource,
				FormatTextHighlight: highlightSource,
			},
		})
	}

	return annotatedWarnings, nil
}
