package version

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Service struct {
	version    string
	buildTime  string
	commitHash string
}

func NewService(
	version string,
	buildTime string,
	commitHash string,
) *Service {
	return &Service{
		version:    version,
		buildTime:  buildTime,
		commitHash: commitHash,
	}
}

func (s *Service) Behave() (models.Version, error) {
	return models.Version{
		LinterVersion:       s.version,
		GoArchFileSupported: fmt.Sprintf("1, %d", models.SupportedVersion),
		BuildTime:           s.buildTime,
		CommitHash:          s.commitHash,
	}, nil
}
