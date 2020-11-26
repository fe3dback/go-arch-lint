package specvalidator

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type SpecValidator struct {
	validators []validator
}

func NewSpecValidator(validators []validator) *SpecValidator {
	return &SpecValidator{
		validators: validators,
	}
}

func (v *SpecValidator) Validate(spec models.ArchSpec) {
	for _, specValidator := range v.validators {
		// todo: return model
		specValidator.Validate(spec)
	}
}
