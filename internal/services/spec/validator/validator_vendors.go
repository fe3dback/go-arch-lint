package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
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

func (v *validatorVendors) Validate(doc arch.Document) []speca.Notice {
	return make([]speca.Notice, 0)
}
