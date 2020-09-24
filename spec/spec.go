package spec

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/fe3dback/go-arch-lint/spec/annotate"
	"github.com/fe3dback/go-arch-lint/spec/archfile"
	specValidator "github.com/fe3dback/go-arch-lint/spec/validator"
	"github.com/goccy/go-yaml"
)

type (
	YamlParseError struct {
		Err      error
		Warnings []YamlAnnotatedWarning
	}

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
)

func newSpec(archFile string, rootDirectory string) (archfile.YamlSpec, YamlParseError) {
	spec := archfile.YamlSpec{}

	sourceCode, err := ioutil.ReadFile(archFile)
	if err != nil {
		return spec, YamlParseError{
			Err:      fmt.Errorf("can`t open '%s': %v", archFile, err),
			Warnings: nil,
		}
	}

	reader := bytes.NewBuffer(sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)
	err = decoder.Decode(&spec)
	if err != nil {
		return spec, YamlParseError{
			Err:      fmt.Errorf("can`t parse yaml in '%s': %v", archFile, err),
			Warnings: nil,
		}
	}

	validator := specValidator.NewArchFileValidator(spec, rootDirectory)
	warnings, err := annotateWarnings(sourceCode, validator.Validate())

	if len(warnings) > 0 {
		return spec, YamlParseError{
			Err:      fmt.Errorf("spec '%s' has warnings", archFile),
			Warnings: warnings,
		}
	}

	return spec, YamlParseError{}
}

func annotateWarnings(sourceCode []byte, warnings []specValidator.Warning) ([]YamlAnnotatedWarning, error) {
	annotatedWarnings := make([]YamlAnnotatedWarning, 0)

	for _, warning := range warnings {
		path := warning.Path

		sourceLine, err := yaml.PathString(path)
		if err != nil {
			return nil, fmt.Errorf("failed check '%s': %v", path, err)
		}

		textSource, err := sourceLine.AnnotateSource(sourceCode, false)
		if err != nil {
			return nil, fmt.Errorf("failed annotate '%s': %v", path, err)
		}

		highlightSource, err := sourceLine.AnnotateSource(sourceCode, true)
		if err != nil {
			return nil, fmt.Errorf("failed annotate '%s': %v", path, err)
		}

		sourceMarker := annotate.ParseSourceError(string(textSource))

		annotatedWarnings = append(annotatedWarnings, YamlAnnotatedWarning{
			Text:   warning.Warning.Error(),
			Path:   warning.Path,
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
