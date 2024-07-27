package validator

import (
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Root struct {
	skipMissuse bool
	validators  []internalValidator
}

func NewRoot(skipMissuse bool, validators ...internalValidator) *Root {
	return &Root{
		skipMissuse: skipMissuse,
		validators:  validators,
	}
}

func (v *Root) Validate(config models.Config) error {
	ctx := &validationContext{
		conf:    config,
		notices: make([]models.Notice, 0, 16),
	}

	for _, validator := range v.validators {
		validator.Validate(ctx)

		if ctx.critical {
			break
		}
	}

	if len(ctx.notices) > 0 {
		return models.NewErrorWithNotices(
			"Config validator find some notices",
			ctx.notices,
		)
	}

	if !v.skipMissuse && len(ctx.missUsage) > 0 {
		return models.NewErrorWithNotices(
			"Config validator find miss usages. You can hide this message by adding '--skip-missuse' flag",
			ctx.missUsage,
		)
	}

	// todo missuse:
	// - allowAnyVendorImports with not empty vendors anywhere
	// - structTags with deps (collision)
	// - deps contain cmp/vendor that defined in common section

	return nil
}
