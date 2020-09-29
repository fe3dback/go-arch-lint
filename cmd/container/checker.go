package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/checker"
)

func (c *Container) ProvideChecker() *checker.Checker {
	resolver := c.provideFilesResolver()
	projectFiles, err := resolver.Resolve()
	if err != nil {
		panic(fmt.Errorf("failed resolve project files: %w", err))
	}

	return checker.NewChecker(
		c.provideProjectRootDirectory(),
		c.provideArch(),
		projectFiles,
	)
}
