package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/models"

	"github.com/fe3dback/go-arch-lint/spec/archfile"
)

func (c *Container) provideArchSpec() *archfile.YamlSpec {
	sourceCode := c.provideArchFileSourceCode()

	archSpec, err := archfile.NewYamlSpec(sourceCode)
	if err != nil {
		panic(models.NewUserSpaceError(fmt.Sprintf("failed provide arch spec: %s", err)))
	}

	return archSpec
}
