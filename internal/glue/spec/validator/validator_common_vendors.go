package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorCommonVendors struct {
	utils *utils
}

func newValidatorCommonVendors(
	utils *utils,
) *validatorCommonVendors {
	return &validatorCommonVendors{
		utils: utils,
	}
}

func (v *validatorCommonVendors) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for _, vendorName := range doc.CommonVendors().List() {
		if err := v.utils.assertKnownVendor(vendorName.Value()); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    vendorName.Reference(),
			})
		}
	}

	return notices
}
