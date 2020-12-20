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
	notices := make([]speca.Notice, 0)

	for _, vendor := range doc.Vendors().Map() {
		err := v.utils.assertVendorImportPathValid(vendor.ImportPath().Value())
		if err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    vendor.ImportPath().Reference(),
			})
		}
	}

	return notices
}
