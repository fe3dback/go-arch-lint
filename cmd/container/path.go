package container

import (
	"github.com/fe3dback/go-arch-lint/path"
)

func (c *Container) providePathResolver() *path.Resolver {
	return path.NewResolver()
}
