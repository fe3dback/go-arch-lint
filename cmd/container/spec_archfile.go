package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/spec/archfile"
)

func (c *Container) provideArchSpec() *archfile.YamlSpec {
	sourceCode := c.provideArchFileSourceCode()

	archSpec, err := archfile.NewYamlSpec(sourceCode)
	if err != nil {
		panic(fmt.Errorf("failed provide arch spec: %w", err))
	}

	return archSpec
}
