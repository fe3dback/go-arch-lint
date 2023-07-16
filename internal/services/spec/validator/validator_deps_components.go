package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
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

func (v *validatorDepsComponents) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for name, rules := range doc.Dependencies().Map() {
		existComponents := make(map[string]bool)

		for _, componentName := range rules.MayDependOn() {
			if _, ok := existComponents[componentName.Value]; ok {
				notices = append(notices, speca.Notice{
					Notice: fmt.Errorf("component '%s' dublicated in '%s' deps", componentName.Value, name),
					Ref:    componentName.Reference,
				})
			}

			if err := v.utils.assertKnownComponent(componentName.Value); err != nil {
				notices = append(notices, speca.Notice{
					Notice: err,
					Ref:    componentName.Reference,
				})
			}

			existComponents[componentName.Value] = true
		}
	}

	return notices
}
