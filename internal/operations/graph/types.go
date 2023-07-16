package graph

import (
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type (
	specAssembler interface {
		Assemble(prj common.Project) (arch.Spec, error)
	}

	projectInfoAssembler interface {
		ProjectInfo(rootDirectory string, archFilePath string) (common.Project, error)
	}
)
