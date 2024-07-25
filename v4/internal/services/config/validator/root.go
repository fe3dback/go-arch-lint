package validator

import (
	"errors"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Root struct {
	validators []internalValidator
}

func NewRoot(validators ...internalValidator) *Root {
	return &Root{
		validators: validators,
	}
}

func (v *Root) Validate(config models.Config) error {
	var resultErr error

	for _, validator := range v.validators {
		resultErr = errors.Join(resultErr, validator.Validate(config))
	}

	return resultErr
}
