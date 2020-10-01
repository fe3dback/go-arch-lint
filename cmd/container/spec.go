package container

import (
	"fmt"
	"io/ioutil"

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
		panic(fmt.Errorf("failed provide arch: %w", err))
	}

	return arch
}

func (c *Container) provideArchFileSourceCode() []byte {
	sourceCode, err := ioutil.ReadFile(
		c.provideArchFilePath(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to provide source code of archfile: %w", err))
	}

	return sourceCode
}
