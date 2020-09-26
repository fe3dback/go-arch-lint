package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/fe3dback/go-arch-lint/checker"
	"github.com/fe3dback/go-arch-lint/files"
	"github.com/fe3dback/go-arch-lint/spec"
	"github.com/fe3dback/go-arch-lint/spec/archfile"
	"github.com/fe3dback/go-arch-lint/spec/validator"
)

// todo:
// - all constructors move to di container
// - split all classes to sep packages, bind with other code wia interfaces
// - unit test's for separated packages

type (
	Container struct {
		archFilePath         string
		projectRootDirectory string
		moduleName           string
	}
)

func newContainer(archFilePath, projectRootDirectory, moduleName string) *Container {
	return &Container{
		archFilePath:         archFilePath,
		projectRootDirectory: projectRootDirectory,
		moduleName:           moduleName,
	}
}

func (c *Container) provideArchFilePath() string {
	return c.archFilePath
}

func (c *Container) provideModuleName() string {
	return c.moduleName
}

func (c *Container) provideProjectRootDirectory() string {
	return c.projectRootDirectory
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

func (c *Container) provideArchSpecValidator() *validator.ArchFileValidator {
	return validator.NewArchFileValidator(
		c.provideArchSpec(),
		c.provideProjectRootDirectory(),
	)
}

func (c *Container) provideArchSpecAnnotatedValidator() *validator.AnnotatedValidator {
	return validator.NewAnnotatedValidator(
		c.provideArchSpecValidator(),
		c.provideArchFileSourceCode(),
	)
}

func (c *Container) provideArchSpec() *archfile.YamlSpec {
	sourceCode := c.provideArchFileSourceCode()

	archSpec, err := archfile.NewYamlSpec(sourceCode)
	if err != nil {
		panic(fmt.Errorf("failed provide arch spec: %w", err))
	}

	return archSpec
}

func (c *Container) provideArch() *spec.Arch {
	arch, err := spec.NewArch(
		c.provideArchSpec(),
		c.provideModuleName(),
		c.provideProjectRootDirectory(),
	)
	if err != nil {
		panic(fmt.Errorf("failed provide arch: %w", err))
	}

	return arch
}

func (c *Container) provideFilesResolver() *files.Resolver {
	arch := c.provideArch()

	excludePaths := make([]string, 0)
	for _, excludeDir := range arch.Exclude {
		excludePaths = append(excludePaths, excludeDir.AbsPath)
	}

	resolver := files.NewResolver(
		settingsProjectDirectory,
		settingsModuleName,
		excludePaths,
		arch.ExcludeFilesMatcher,
	)

	return resolver
}

func (c *Container) provideChecker() *checker.Checker {
	resolver := c.provideFilesResolver()
	projectFiles, err := resolver.Resolve()
	if err != nil {
		panic(fmt.Errorf("failed resolve project files: %w", err))
	}

	return checker.NewChecker(
		settingsProjectDirectory,
		c.provideArch(),
		projectFiles,
	)
}
