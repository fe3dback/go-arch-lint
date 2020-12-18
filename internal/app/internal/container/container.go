package container

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Container struct {
	version    string
	buildTime  string
	commitHash string

	flags models.FlagsRoot
}

func NewContainer(
	version string,
	buildTime string,
	commitHash string,
) *Container {
	return &Container{
		version:    version,
		buildTime:  buildTime,
		commitHash: commitHash,
	}
}
