package selfInspect

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	projectInfoAssembler interface {
		ProjectInfo(rootDirectory string, archFilePath string) (models.ProjectInfo, error)
	}
)
