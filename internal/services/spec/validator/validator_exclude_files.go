package validator

import (
	"fmt"
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type validatorExcludeFiles struct {
}

func newValidatorExcludeFiles() *validatorExcludeFiles {
	return &validatorExcludeFiles{}
}

func (v *validatorExcludeFiles) Validate(doc spec.Document) []arch.Notice {
	notices := make([]arch.Notice, 0)

	for index, regExp := range doc.ExcludedFilesRegExp() {
		if _, err := regexp.Compile(regExp.Value); err != nil {
			notices = append(notices, arch.Notice{
				Notice: fmt.Errorf("invalid regexp '%s' at %d: %v", regExp.Value, index, err),
				Ref:    regExp.Reference,
			})
		}
	}

	return notices
}
