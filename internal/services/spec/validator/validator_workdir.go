package validator

import (
	"fmt"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type validatorWorkDir struct {
	utils *utils
}

func newValidatorWorkDir(utils *utils) *validatorWorkDir {
	return &validatorWorkDir{
		utils: utils,
	}
}

func (v *validatorWorkDir) Validate(doc spec.Document) []arch.Notice {
	notices := make([]arch.Notice, 0)

	// workdir comes from archfile in unix form, projectDir is os specific,
	// filepath.Join already cleans result, extra path.Clean does nothing on windows
	absPath := filepath.Join(v.utils.projectDir, filepath.FromSlash(doc.WorkingDirectory().Value))

	err := v.utils.assertDirectoriesValid(absPath)
	if err != nil {
		notices = append(notices, arch.Notice{
			Notice: fmt.Errorf("invalid workdir '%s' (%s), directory not exist",
				doc.WorkingDirectory().Value,
				absPath,
			),
			Ref: doc.WorkingDirectory().Reference,
		})
	}

	return notices
}
