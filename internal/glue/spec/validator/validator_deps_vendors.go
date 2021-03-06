package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorDepsVendors struct {
	utils *utils
}

func newValidatorDepsVendors(
	utils *utils,
) *validatorDepsVendors {
	return &validatorDepsVendors{
		utils: utils,
	}
}

func (v *validatorDepsVendors) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for name, rules := range doc.Dependencies().Map() {
		existVendors := make(map[string]bool)

		for _, vendorName := range rules.CanUse() {
			if _, ok := existVendors[vendorName.Value()]; ok {
				notices = append(notices, speca.Notice{
					Notice: fmt.Errorf("vendor '%s' dublicated in '%s' deps", vendorName.Value(), name),
					Ref:    vendorName.Reference(),
				})
			}

			if err := v.utils.assertKnownVendor(vendorName.Value()); err != nil {
				notices = append(notices, speca.Notice{
					Notice: err,
					Ref:    vendorName.Reference(),
				})
			}

			existVendors[vendorName.Value()] = true
		}
	}

	return notices
}
