package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorCommonComponents struct {
	refResolver YamlSourceCodeReferenceResolver
	utils       *utils
}

func newValidatorCommonComponents(
	refResolver YamlSourceCodeReferenceResolver,
	utils *utils,
) *validatorCommonComponents {
	return &validatorCommonComponents{
		refResolver: refResolver,
		utils:       utils,
	}
}

func (v *validatorCommonComponents) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for index, componentName := range doc.CommonComponents {
		if err := v.utils.assertKnownComponent(componentName); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    v.refResolver.Resolve(fmt.Sprintf("$.commonComponents[%d]", index)),
			})
		}
	}

	return notices
}
