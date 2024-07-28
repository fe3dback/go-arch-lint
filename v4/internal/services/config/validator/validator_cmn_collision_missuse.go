package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type VendorsCommonCollisionMissUseValidator struct{}

func NewCommonCollisionMissUseValidator() *VendorsCommonCollisionMissUseValidator {
	return &VendorsCommonCollisionMissUseValidator{}
}

func (c *VendorsCommonCollisionMissUseValidator) Validate(ctx *validationContext) {
	ctx.conf.Dependencies.Map.Each(func(name models.ComponentName, rules models.ConfigComponentDependencies, _ models.Reference) {
		for _, anotherComponentName := range rules.MayDependOn {
			if !ctx.conf.CommonComponents.Contains(anotherComponentName) {
				continue
			}

			ctx.AddMissUse(
				fmt.Sprintf("redundant declaration: component '%s' may depend on '%s', but '%s' is already \"common component\".",
					name,
					anotherComponentName.Value,
					anotherComponentName.Value,
				),
				anotherComponentName.Ref,
			)
		}

		for _, anotherVendorName := range rules.CanUse {
			if !ctx.conf.CommonVendors.Contains(anotherVendorName) {
				continue
			}

			ctx.AddMissUse(
				fmt.Sprintf("redundant declaration: component '%s' can use '%s', but '%s' is already \"common vendor\".",
					name,
					anotherVendorName.Value,
					anotherVendorName.Value,
				),
				anotherVendorName.Ref,
			)
		}
	})
}
