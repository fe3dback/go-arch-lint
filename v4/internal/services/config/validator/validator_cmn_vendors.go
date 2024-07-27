package validator

import "fmt"

type CommonVendorsValidator struct{}

func NewCommonVendorsValidator() *CommonVendorsValidator {
	return &CommonVendorsValidator{}
}

func (c *CommonVendorsValidator) Validate(ctx *validationContext) {
	for _, name := range ctx.conf.CommonVendors {
		if !ctx.IsKnownVendor(name.Value) {
			ctx.AddNotice(
				fmt.Sprintf("Common vendor '%s' is not known", name.Value),
				name.Ref,
			)
		}
	}
}
