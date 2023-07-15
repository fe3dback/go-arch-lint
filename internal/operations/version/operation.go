package version

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Operation struct {
	version    string
	buildTime  string
	commitHash string
}

func NewOperation(
	version string,
	buildTime string,
	commitHash string,
) *Operation {
	return &Operation{
		version:    version,
		buildTime:  buildTime,
		commitHash: commitHash,
	}
}

func (s *Operation) Behave() (models.CmdVersionOut, error) {
	return models.CmdVersionOut{
		LinterVersion:       s.version,
		GoArchFileSupported: fmt.Sprintf("1 .. %d", models.SupportedVersion),
		BuildTime:           s.buildTime,
		CommitHash:          s.commitHash,
	}, nil
}
