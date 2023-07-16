package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type validatorVersion struct{}

func newValidatorVersion() *validatorVersion {
	return &validatorVersion{}
}

func (v *validatorVersion) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	if doc.Version().Value <= models.SupportedVersionMax && doc.Version().Value >= models.SupportedVersionMin {
		return notices
	}

	notices = append(notices, speca.Notice{
		Notice: fmt.Errorf("version '%d' is not supported, supported: [%d-%d]",
			doc.Version().Value,
			models.SupportedVersionMin,
			models.SupportedVersionMax,
		),
		Ref: doc.Version().Reference,
	})

	return notices
}
