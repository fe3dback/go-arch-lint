package assembler

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	projectInfoFetcher interface {
		Fetch() (models.ProjectInfo, error)
	}

	pathHelper interface {
		MatchProjectFiles(somePath any, queryType models.FileMatchQueryType) ([]models.FileRef, error)
	}
)
