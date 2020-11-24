package container

import (
	"github.com/fe3dback/go-arch-lint/internal/services/version"
)

func (c *Container) provideVersionService() *version.Service {
	return version.NewService(
		c.version,
		c.buildTime,
		c.commitHash,
	)
}
