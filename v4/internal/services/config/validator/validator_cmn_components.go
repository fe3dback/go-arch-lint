package validator

import "fmt"

type CommonComponentsValidator struct{}

func NewCommonComponentsValidator() *CommonComponentsValidator {
	return &CommonComponentsValidator{}
}

func (c *CommonComponentsValidator) Validate(ctx *validationContext) {
	for _, name := range ctx.conf.CommonComponents {
		if !ctx.IsKnownComponent(name.Value) {
			ctx.AddNotice(
				fmt.Sprintf("Common component '%s' is not known", name.Value),
				name.Ref,
			)
		}
	}
}
