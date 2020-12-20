package validator

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorCommonComponents struct {
	utils *utils
}

func newValidatorCommonComponents(
	utils *utils,
) *validatorCommonComponents {
	return &validatorCommonComponents{
		utils: utils,
	}
}

func (v *validatorCommonComponents) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for _, componentName := range doc.CommonComponents().List() {
		if err := v.utils.assertKnownComponent(componentName.Value()); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    componentName.Reference(),
			})
		}
	}

	return notices
}
