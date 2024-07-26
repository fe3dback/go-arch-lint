package validator

import (
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Root struct {
	validators []internalValidator
}

func NewRoot(validators ...internalValidator) *Root {
	return &Root{
		validators: validators,
	}
}

func (v *Root) Validate(config models.Config) error {
	ctx := &validationContext{
		conf:    config,
		notices: make([]models.Notice, 0, 16),
	}

	for _, validator := range v.validators {
		validator.Validate(ctx)
	}

	if len(ctx.notices) > 0 {
		return models.NewErrorWithNotices(
			"Config validation found some notices",
			ctx.notices,
		)
	}

	return nil
}
