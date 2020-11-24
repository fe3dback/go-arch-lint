package version

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

const (
	goArchFileSupported = "1"
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

func (p *Service) Behave() (models.Version, error) {
	return models.Version{
		LinterVersion:       p.version,
		GoArchFileSupported: goArchFileSupported,
		BuildTime:           p.buildTime,
		CommitHash:          p.commitHash,
	}, nil
}
