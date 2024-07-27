package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type DepsValidator struct {
	pathHelper pathHelper
}

func NewDepsValidator(
	pathHelper pathHelper,
) *DepsValidator {
	return &DepsValidator{
		pathHelper: pathHelper,
	}
}

func (c *DepsValidator) Validate(ctx *validationContext) {
	ctx.conf.Dependencies.Map.Each(func(name models.ComponentName, rules models.ConfigComponentDependencies, reference models.Reference) {
		if !ctx.IsKnownComponent(name) {
			ctx.AddNotice(
				fmt.Sprintf("Component '%s' in dependencies is not defined", name),
				reference,
			)
			return
		}

		if rules.AnyProjectDeps.Value && len(rules.MayDependOn) > 0 {
			ctx.AddMissUse(
				fmt.Sprintf("In component '%s': rule 'anyProjectDeps' used with not empty 'MayDependOn' list", name),
				rules.AnyProjectDeps.Ref,
			)
			return
		}

		if rules.AnyVendorDeps.Value && len(rules.CanUse) > 0 {
			ctx.AddMissUse(
				fmt.Sprintf("In component '%s': rule 'anyVendorDeps' used with not empty 'CanUse' list", name),
				rules.AnyVendorDeps.Ref,
			)
			return
		}

		if len(rules.MayDependOn) == 0 && len(rules.CanUse) == 0 {
			if rules.AnyProjectDeps.Value {
				return
			}

			if rules.AnyVendorDeps.Value {
				return
			}

			ctx.AddNotice(
				fmt.Sprintf("In component '%s': no rules is defined (require at least one of [AnyProjectDeps, AnyVendorDeps, MayDependOn, CanUse])", name),
				reference,
			)
			return
		}
	})
}
