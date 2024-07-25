package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type CommonComponentsValidator struct{}

func NewCommonComponentsValidator() *CommonComponentsValidator {
	return &CommonComponentsValidator{}
}

func (c *CommonComponentsValidator) Validate(conf models.Config) error {
	return models.NewReferencedError(
		fmt.Errorf("version should be as string"),
		conf.Version.Ref,
	)
}
