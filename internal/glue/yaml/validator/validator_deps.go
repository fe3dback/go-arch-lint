package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorDeps struct {
	refResolver YamlSourceCodeReferenceResolver
	utils       *utils
}

func newValidatorDeps(
	refResolver YamlSourceCodeReferenceResolver,
	utils *utils,
) *validatorDeps {
	return &validatorDeps{
		refResolver: refResolver,
		utils:       utils,
	}
}

func (v *validatorDeps) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for name, rules := range doc.Dependencies {
		refPath := fmt.Sprintf("$.deps.%s", name)

		if err := v.utils.assertKnownComponent(name); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    v.refResolver.Resolve(refPath),
			})
		}

		if len(rules.MayDependOn) == 0 && len(rules.CanUse) == 0 {
			if rules.AnyProjectDeps {
				continue
			}

			if rules.AnyVendorDeps {
				continue
			}

			notices = append(notices, speca.Notice{
				Notice: fmt.Errorf("should have ref in 'mayDependOn' or at least one flag of ['anyProjectDeps', 'anyVendorDeps']"),
				Ref:    v.refResolver.Resolve(refPath),
			})
		}
	}

	return notices
}
