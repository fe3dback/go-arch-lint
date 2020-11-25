package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/models"

	"github.com/fe3dback/go-arch-lint/spec"
)

func (c *Container) provideArch() *spec.Arch {
	arch, err := spec.NewArch(
		c.providePathResolver(),
		c.provideArchSpec(),
		c.provideModuleName(),
		c.provideProjectRootDirectory(),
	)
	if err != nil {
		panic(models.NewUserSpaceError(fmt.Sprintf("failed provide arch: %s", err)))
	}

	return arch
}

func (c *Container) provideArchFileSourceCode() []byte {
}
