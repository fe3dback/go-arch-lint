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
		SourceCode []byte `json:"-"`
	}

	AnnotatedValidator struct {
		innerValidator Validator
		warningParser  AnnotatedWarningParser
		sourceCode     []byte
		useColors      bool
	}
)

func NewAnnotatedValidator(
	innerValidator Validator,
	annotatedWarningParser AnnotatedWarningParser,
	sourceCode []byte,
	useColors bool,
) *AnnotatedValidator {
	return &AnnotatedValidator{
		innerValidator: innerValidator,
		warningParser:  annotatedWarningParser,
		sourceCode:     sourceCode,
		useColors:      useColors,
	}
}

func (av *AnnotatedValidator) Validate() ([]YamlAnnotatedWarning, error) {
	return av.annotateWarnings(
		av.innerValidator.Validate(),
	)
}

func (av *AnnotatedValidator) annotateWarnings(
	warnings []models.ArchFileSyntaxWarning,
) (annotatedWarnings []YamlAnnotatedWarning, executionErr error) {
	annotatedWarnings = make([]YamlAnnotatedWarning, 0)

	for _, warning := range warnings {
		path := warning.Path()
		actualError := warning.Warning()

		code, line, pos, err := av.annotateWarning(path)
		if err != nil {
			actualError = fmt.Errorf("%v (+annotation: %v)", actualError, err)
		}

		if actualError == nil {
			actualError = fmt.Errorf("empty err for path '%s'", path)
		}

		annotatedWarnings = append(annotatedWarnings, YamlAnnotatedWarning{
			Text:       actualError.Error(),
			Path:       path,
			Line:       line,
			Offset:     pos,
			SourceCode: code,
		})
	}

	return annotatedWarnings, nil
}

func (av *AnnotatedValidator) annotateWarning(path string) (
	codeWidget []byte, line int, pos int, executionErr error,
) {
	defer func() {
		if err := recover(); err != nil {
			executionErr = fmt.Errorf("failed annotate path '%s' failed: %s", path, err)
			return
		}
	}()

	sourceLine, err := yaml.PathString(path)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed check '%s': %v", path, err)
	}

	rawSource, err := sourceLine.AnnotateSource(av.sourceCode, false)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed annotate source in '%s': %v", path, err)
	}

	line, pos, err = av.warningParser.Parse(string(rawSource))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed parse source warning text for path '%s': %v", path, err)
	}

	textSource, err := sourceLine.AnnotateSource(av.sourceCode, av.useColors)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed annotate source in '%s': %v", path, err)
	}

	return textSource, line, pos, nil
}
