package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorComponents struct {
	refResolver YamlSourceCodeReferenceResolver
	utils       *utils
}

func newValidatorComponents(
	refResolver YamlSourceCodeReferenceResolver,
	utils *utils,
) *validatorComponents {
	return &validatorComponents{
		refResolver: refResolver,
		utils:       utils,
	}
}

func (v *validatorComponents) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	if len(doc.Components) == 0 {
		notices = append(notices, speca.Notice{
			Notice: fmt.Errorf("at least one component should by defined"),
			Ref:    v.refResolver.Resolve("$.components"),
		})
	}

	for name, component := range doc.Components {
		if err := v.utils.assertPathValid(component.LocalPath); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    v.refResolver.Resolve(fmt.Sprintf("$.components.%s.in", name)),
			})
		}
	}

	return notices
}
