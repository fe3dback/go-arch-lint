package validator

import (
	"fmt"
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorExcludeFiles struct {
	refResolver YamlSourceCodeReferenceResolver
}

func newValidatorExcludeFiles(refResolver YamlSourceCodeReferenceResolver) *validatorExcludeFiles {
	return &validatorExcludeFiles{refResolver: refResolver}
}

func (v *validatorExcludeFiles) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for index, regExp := range doc.ExcludeFilesRegExp {
		if _, err := regexp.Compile(regExp); err != nil {
			notices = append(notices, speca.Notice{
				Notice: fmt.Errorf("invalid regexp '%s' at %d: %v", regExp, index, err),
				Ref:    v.refResolver.Resolve(fmt.Sprintf("$.excludeFiles[%d]", index)),
			})
		}
	}

	return notices
}
