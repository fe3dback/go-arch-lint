package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorComponents struct {
	utils *utils
}

func newValidatorComponents(
	utils *utils,
) *validatorComponents {
	return &validatorComponents{
		utils: utils,
	}
}

func (v *validatorComponents) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	if len(doc.Components().Map()) == 0 {
		notices = append(notices, speca.Notice{
			Notice: fmt.Errorf("at least one component should by defined"),
			Ref:    doc.Components().Reference(),
		})
	}

	for _, component := range doc.Components().Map() {
		localPath := component.LocalPath()
		if err := v.utils.assertPathValid(localPath.Value()); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    localPath.Reference(),
			})
		}
	}

	return notices
}
