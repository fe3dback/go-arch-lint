package main

import (
	"os"

	"github.com/fe3dback/go-arch-lint/checker"

	"github.com/fe3dback/go-arch-lint/files"
	"github.com/fe3dback/go-arch-lint/spec"
)

func main() {
	flags, err := flagsParse()
	if err != nil {
		panic(err)
	}

	arch, err := spec.NewArch(
		flags.pathArchFile,
		flags.moduleName,
		flags.pathProjectDirectory,
	)
	if err != nil {
		panic(err)
	}

	resolver := createResolver(arch, flags)
	projectFiles, err := resolver.Resolve()
	if err != nil {
		panic(err)
	}

	archChecker := checker.NewChecker(
		flags.pathProjectDirectory,
		arch,
		projectFiles,
		os.Stdout,
	)
	archChecker.Check()
}

func createResolver(arch *spec.Arch, flags flags) *files.Resolver {
	excludePaths := make([]string, 0)
	for _, excludeDir := range arch.Exclude {
		excludePaths = append(excludePaths, excludeDir.AbsPath)
	}

	resolver := files.NewResolver(
		flags.pathProjectDirectory,
		flags.moduleName,
		excludePaths,
		arch.ExcludeFilesMatcher,
	)

	return resolver
}
