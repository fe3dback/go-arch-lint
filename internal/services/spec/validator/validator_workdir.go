package validator

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorWorkDir struct {
	utils *utils
}

func newValidatorWorkDir(utils *utils) *validatorWorkDir {
	return &validatorWorkDir{
		utils: utils,
	}
}

func (v *validatorWorkDir) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	absPath := fmt.Sprintf("%s/%s", v.utils.rootDirectory, doc.WorkingDirectory().Value())
	absPath = path.Clean(absPath)

	err := v.utils.assertDirectoriesValid(absPath)
	if err != nil {
		notices = append(notices, speca.Notice{
			Notice: fmt.Errorf("invalid workdir '%s' (%s), directory not exist",
				doc.WorkingDirectory().Value(),
				absPath,
			),
			Ref: doc.WorkingDirectory().Reference(),
		})
	}

	return notices
}
