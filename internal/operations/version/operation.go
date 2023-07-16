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

func (o *Operation) Behave() (models.CmdVersionOut, error) {
	return models.CmdVersionOut{
		LinterVersion:       o.version,
		GoArchFileSupported: fmt.Sprintf("%d .. %d", models.SupportedVersionMin, models.SupportedVersionMax),
		BuildTime:           o.buildTime,
		CommitHash:          o.commitHash,
	}, nil
}
