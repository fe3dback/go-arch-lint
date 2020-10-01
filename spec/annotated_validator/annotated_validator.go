package annotated_validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/models"
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
		innerValidator Validator
		warningParser  AnnotatedWarningParser
		sourceCode     []byte
	}
)

func NewAnnotatedValidator(
	innerValidator Validator,
	annotatedWarningParser AnnotatedWarningParser,
	sourceCode []byte,
) *AnnotatedValidator {
	return &AnnotatedValidator{
		innerValidator: innerValidator,
		warningParser:  annotatedWarningParser,
		sourceCode:     sourceCode,
	}
}

func (av *AnnotatedValidator) Validate() ([]YamlAnnotatedWarning, error) {
	return av.annotateWarnings(
		av.innerValidator.Validate(),
	)
}

func (av *AnnotatedValidator) annotateWarnings(warnings []models.ArchFileSyntaxWarning) ([]YamlAnnotatedWarning, error) {
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

		line, pos, err := av.warningParser.Parse(string(textSource))
		if err != nil {
			return nil, fmt.Errorf("failed parse source warning text for path '%s': %v", path, err)
		}

		annotatedWarnings = append(annotatedWarnings, YamlAnnotatedWarning{
			Text:   warning.Warning().Error(),
			Path:   path,
			Line:   line,
			Offset: pos,
			SourceCode: &YamlAnnotatedWarningSource{
				FormatText:          textSource,
				FormatTextHighlight: highlightSource,
			},
		})
	}

	return annotatedWarnings, nil
}
