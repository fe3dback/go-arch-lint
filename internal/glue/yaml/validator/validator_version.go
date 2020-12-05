package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

// todo: some version management
const supportedVersion = 1

type validatorVersion struct {
	refResolver YamlSourceCodeReferenceResolver
}

func newValidatorVersion(refResolver YamlSourceCodeReferenceResolver) *validatorVersion {
	return &validatorVersion{refResolver: refResolver}
}

func (v *validatorVersion) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	if doc.Version <= supportedVersion && doc.Version > 0 {
		return notices
	}

	notices = append(notices, speca.Notice{
		Notice: fmt.Errorf("version '%d' is not supported, supported: [%d]",
			doc.Version,
			supportedVersion,
		),
		Ref: v.refResolver.Resolve("$.version"),
	})

	return notices
}
