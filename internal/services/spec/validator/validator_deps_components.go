package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type validatorDepsComponents struct {
	utils *utils
}

func newValidatorDepsComponents(
	utils *utils,
) *validatorDepsComponents {
	return &validatorDepsComponents{
		utils: utils,
	}
}

func (v *validatorDepsComponents) Validate(doc spec.Document) []arch.Notice {
	notices := make([]arch.Notice, 0)

	for name, rule := range doc.Dependencies() {
		existComponents := make(map[string]bool)

		for _, componentName := range rule.Value.MayDependOn() {
			if _, ok := existComponents[componentName.Value]; ok {
				notices = append(notices, arch.Notice{
					Notice: fmt.Errorf("component '%s' dublicated in '%s' deps", componentName.Value, name),
					Ref:    componentName.Reference,
				})
			}

			if err := v.utils.assertKnownComponent(componentName.Value); err != nil {
				notices = append(notices, arch.Notice{
					Notice: err,
					Ref:    componentName.Reference,
				})
			}

			existComponents[componentName.Value] = true
		}
	}

	return notices
}
