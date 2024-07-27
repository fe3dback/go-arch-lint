package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type DepsComponentsValidator struct {
	pathHelper pathHelper
}

func NewDepsComponentsValidator(
	pathHelper pathHelper,
) *DepsComponentsValidator {
	return &DepsComponentsValidator{
		pathHelper: pathHelper,
	}
}

func (c *DepsComponentsValidator) Validate(ctx *validationContext) {
	ctx.conf.Dependencies.Map.Each(func(name models.ComponentName, rules models.ConfigComponentDependencies, reference models.Reference) {
		existComponents := make(map[models.ComponentName]any)

		for _, anotherComponentName := range rules.MayDependOn {
			// check multiple usage
			if _, ok := existComponents[anotherComponentName.Value]; ok {
				ctx.AddNotice(
					fmt.Sprintf("Component '%s' dublicated in '%s' deps", anotherComponentName.Value, name),
					anotherComponentName.Ref,
				)
			}

			existComponents[anotherComponentName.Value] = struct{}{}

			// check is known
			if !ctx.IsKnownComponent(anotherComponentName.Value) {
				ctx.AddNotice(
					fmt.Sprintf("Component '%s' is not defined", anotherComponentName.Value),
					anotherComponentName.Ref,
				)
			}
		}
	})
}
