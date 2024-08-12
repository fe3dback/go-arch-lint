package assembler

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	projectInfoFetcher interface {
		Fetch() (models.ProjectInfo, error)
	}

	pathHelper interface {
		FindProjectFiles(query models.FileQuery) ([]models.FileDescriptor, error)
	}
)
