package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type DepsVendorsValidator struct {
	pathHelper pathHelper
}

func NewDepsVendorsValidator(
	pathHelper pathHelper,
) *DepsVendorsValidator {
	return &DepsVendorsValidator{
		pathHelper: pathHelper,
	}
}

func (c *DepsVendorsValidator) Validate(ctx *validationContext) {
	ctx.conf.Dependencies.Map.Each(func(name models.ComponentName, rules models.ConfigComponentDependencies, reference models.Reference) {
		existVendors := make(map[models.VendorName]any)

		for _, vendorName := range rules.CanUse {
			// check multiple usage
			if _, ok := existVendors[vendorName.Value]; ok {
				ctx.AddNotice(
					fmt.Sprintf("Vendor '%s' dublicated in '%s' deps", vendorName.Value, name),
					vendorName.Ref,
				)
			}

			existVendors[vendorName.Value] = struct{}{}

			// check is known
			if !ctx.IsKnownVendor(vendorName.Value) {
				ctx.AddNotice(
					fmt.Sprintf("Vendor '%s' is not defined", vendorName.Value),
					vendorName.Ref,
				)
			}
		}
	})
}
