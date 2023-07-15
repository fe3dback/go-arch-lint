package mapping

import (
	"context"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	specAssembler interface {
		Assemble(prj common.Project) (speca.Spec, error)
	}

	projectFilesResolver interface {
		ProjectFiles(ctx context.Context, spec speca.Spec) ([]models.FileHold, error)
	}

	projectInfoAssembler interface {
		ProjectInfo(rootDirectory string, archFilePath string) (common.Project, error)
	}
)
