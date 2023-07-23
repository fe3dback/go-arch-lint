package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type validatorVendors struct {
	utils *utils
}

func newValidatorVendors(
	utils *utils,
) *validatorVendors {
	return &validatorVendors{
		utils: utils,
	}
}

func (v *validatorVendors) Validate(_ spec.Document) []arch.Notice {
	return make([]arch.Notice, 0)
}
