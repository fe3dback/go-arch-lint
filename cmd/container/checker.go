package container

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/checker"
	"github.com/fe3dback/go-arch-lint/models"
)

func (c *Container) ProvideChecker() *checker.Checker {
	resolver := c.provideFilesResolver()
	projectFiles, err := resolver.Resolve()
	if err != nil {
		panic(models.NewUserSpaceError(fmt.Sprintf("failed resolve project files: %s", err)))
	}

	return checker.NewChecker(
		c.provideProjectRootDirectory(),
		c.provideArch(),
		projectFiles,
	)
}
