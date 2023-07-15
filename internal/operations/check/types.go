package check

import (
	"context"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	projectInfoAssembler interface {
		ProjectInfo(rootDirectory string, archFilePath string) (common.Project, error)
	}

	specAssembler interface {
		Assemble(prj common.Project) (speca.Spec, error)
	}

	referenceRender interface {
		SourceCode(ref models.CodeReference, highlight bool) []byte
	}

	specChecker interface {
		Check(ctx context.Context, spec speca.Spec) (models.CheckResult, error)
	}
)
