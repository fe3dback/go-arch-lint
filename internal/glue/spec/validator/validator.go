package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Validator struct {
	validators []validator
}

func NewValidator() *Validator {
	return &Validator{
		validators: []validator{
			// todo: add validators
		},
	}
}

func (v *Validator) Validate(spec speca.Spec) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for _, specValidator := range v.validators {
		notices = append(notices, specValidator.Validate(spec)...)
	}

	return notices
}
