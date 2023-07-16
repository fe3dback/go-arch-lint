package check

import (
	"context"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type (
	projectInfoAssembler interface {
		ProjectInfo(rootDirectory string, archFilePath string) (common.Project, error)
	}

	specAssembler interface {
		Assemble(prj common.Project) (arch.Spec, error)
	}

	referenceRender interface {
		SourceCode(ref common.Reference, highlight bool, showPointer bool) []byte
	}

	specChecker interface {
		Check(ctx context.Context, spec arch.Spec) (models.CheckResult, error)
	}
)
