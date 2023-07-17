package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
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

func (v *validatorCommonVendors) Validate(doc spec.Document) []arch.Notice {
	notices := make([]arch.Notice, 0)

	for _, vendorName := range doc.CommonVendors() {
		if err := v.utils.assertKnownVendor(vendorName.Value); err != nil {
			notices = append(notices, arch.Notice{
				Notice: err,
				Ref:    vendorName.Reference,
			})
		}
	}

	return notices
}
