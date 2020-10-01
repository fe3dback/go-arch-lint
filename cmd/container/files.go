package container

import (
	"github.com/fe3dback/go-arch-lint/files"
)

func (c *Container) provideFilesResolver() *files.Resolver {
	arch := c.provideArch()

	excludePaths := make([]string, 0)
	for _, excludeDir := range arch.Exclude {
		excludePaths = append(excludePaths, excludeDir.AbsPath)
	}

	resolver := files.NewResolver(
		c.provideProjectRootDirectory(),
		c.provideModuleName(),
		excludePaths,
		arch.ExcludeFilesMatcher,
	)

	return resolver
}
