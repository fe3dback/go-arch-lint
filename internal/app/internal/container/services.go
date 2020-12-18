package container

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/services/check"
	"github.com/fe3dback/go-arch-lint/internal/services/mapping"
	"github.com/fe3dback/go-arch-lint/internal/services/version"
)

func (c *Container) provideVersionService() *version.Service {
	return version.NewService(
		c.version,
		c.buildTime,
		c.commitHash,
	)
}

func (c *Container) provideCheckService(input models.FlagsCheck) *check.Service {
	return check.NewService(
		c.provideSpecAssembler(
			input.Project.Directory,
			input.Project.ModuleName,
			input.Project.GoArchFilePath,
		),
		c.provideSpecChecker(),
		c.provideReferenceRender(),
		c.flags.UseColors,
	)
}

func (c *Container) provideMappingService(input models.FlagsMapping) *mapping.Service {
	return mapping.NewService(
		c.provideSpecAssembler(
			input.Project.Directory,
			input.Project.ModuleName,
			input.Project.GoArchFilePath,
		),
		c.provideProjectFilesResolver(),
	)
}
