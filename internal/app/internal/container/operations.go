package container

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/operations/check"
	"github.com/fe3dback/go-arch-lint/internal/operations/mapping"
	"github.com/fe3dback/go-arch-lint/internal/operations/schema"
	"github.com/fe3dback/go-arch-lint/internal/operations/selfInspect"
	"github.com/fe3dback/go-arch-lint/internal/operations/version"
)

func (c *Container) provideOperationVersion() *version.Operation {
	return version.NewOperation(
		c.version,
		c.buildTime,
		c.commitHash,
	)
}

func (c *Container) provideOperationSelfInspect(input models.FlagsSelfInspect) *selfInspect.Operation {
	return selfInspect.NewOperation(
		c.provideSpecAssembler(
			input.Project.Directory,
			input.Project.ModuleName,
			input.Project.GoArchFilePath,
		),
		c.version,
	)
}

func (c *Container) provideOperationCheck(input models.FlagsCheck) *check.Operation {
	return check.NewOperation(
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

func (c *Container) provideOperationMapping(input models.FlagsMapping) *mapping.Operation {
	return mapping.NewOperation(
		c.provideSpecAssembler(
			input.Project.Directory,
			input.Project.ModuleName,
			input.Project.GoArchFilePath,
		),
		c.provideProjectFilesResolver(),
	)
}

func (c *Container) provideOperationSchema() *schema.Operation {
	return schema.NewOperation()
}
