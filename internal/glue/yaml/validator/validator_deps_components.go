package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorDepsComponents struct {
	refResolver YamlSourceCodeReferenceResolver
	utils       *utils
}

func newValidatorDepsComponents(
	refResolver YamlSourceCodeReferenceResolver,
	utils *utils,
) *validatorDepsComponents {
	return &validatorDepsComponents{
		refResolver: refResolver,
		utils:       utils,
	}
}

func (v *validatorDepsComponents) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for name, rules := range doc.Dependencies {
		existComponents := make(map[string]bool)

		for index, componentName := range rules.MayDependOn {
			refPath := fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, index)

			if _, ok := existComponents[componentName]; ok {
				notices = append(notices, speca.Notice{
					Notice: fmt.Errorf("component '%s' dublicated in '%s' deps", componentName, name),
					Ref:    v.refResolver.Resolve(refPath),
				})
			}

			if err := v.utils.assertKnownComponent(componentName); err != nil {
				notices = append(notices, speca.Notice{
					Notice: err,
					Ref:    v.refResolver.Resolve(refPath),
				})
			}

			existComponents[componentName] = true
		}
	}

	return notices
}
