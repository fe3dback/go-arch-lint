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
		for _, vendorIn := range vendor.ImportPaths() {
			v.validateImportPath(&notices, vendorIn)
		}
	}

	return notices
}

func (v *validatorVendors) validateImportPath(notices *[]speca.Notice, vendorIn speca.ReferableString) {
	err := v.utils.assertVendorImportPathValid(vendorIn.Value())
	if err != nil {
		*notices = append(*notices, speca.Notice{
			Notice: err,
			Ref:    vendorIn.Reference(),
		})
	}
}
