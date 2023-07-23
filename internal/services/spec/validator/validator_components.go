package validator

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
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

func (v *validatorComponents) Validate(doc spec.Document) []arch.Notice {
	notices := make([]arch.Notice, 0)

	if len(doc.Components()) == 0 {
		notices = append(notices, arch.Notice{
			Notice: fmt.Errorf("at least one component should by defined"),
			Ref:    doc.Version().Reference,
		})
	}

	for _, component := range doc.Components() {
		for _, componentIn := range component.Value.RelativePaths() {
			localPath := path.Clean(fmt.Sprintf("%s/%s",
				doc.WorkingDirectory().Value,
				string(componentIn),
			))

			if err := v.utils.assertGlobPathValid(localPath); err != nil {
				notices = append(notices, arch.Notice{
					Notice: err,
					Ref:    component.Reference,
				})
			}
		}
	}

	return notices
}
