package graph

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	ProjectInfoAssembler interface {
		ProjectInfo(rootDirectory string, archFilePath string) (models.ProjectInfo, error)
	}
)
