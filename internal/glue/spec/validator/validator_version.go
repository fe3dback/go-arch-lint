package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

// todo: some version management
const supportedVersion = 1

type validatorVersion struct{}

func newValidatorVersion() *validatorVersion {
	return &validatorVersion{}
}

func (v *validatorVersion) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	if doc.Version().Value() <= supportedVersion && doc.Version().Value() > 0 {
		return notices
	}

	notices = append(notices, speca.Notice{
		Notice: fmt.Errorf("version '%d' is not supported, supported: [%d]",
			doc.Version().Value(),
			supportedVersion,
		),
		Ref: doc.Version().Reference(),
	})

	return notices
}
