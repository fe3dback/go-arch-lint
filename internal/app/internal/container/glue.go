package container

import (
	"fmt"
	"io/ioutil"

	"github.com/fe3dback/go-arch-lint/internal/glue/checker"
	"github.com/fe3dback/go-arch-lint/internal/glue/code"
	"github.com/fe3dback/go-arch-lint/internal/glue/path"
	"github.com/fe3dback/go-arch-lint/internal/glue/project"
	specassembler "github.com/fe3dback/go-arch-lint/internal/glue/spec/assembler"
	specvalidator "github.com/fe3dback/go-arch-lint/internal/glue/spec/validator"
	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/reference"
	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/validator"
)

func (c *Container) provideSpecAssembler(projectDir, moduleName, archFilePath string) *specassembler.Assembler {
	return specassembler.NewAssembler(
		c.provideYamlSpecProvider(projectDir, archFilePath),
		c.providePathResolver(),
		c.provideSourceCodeReferenceResolver(archFilePath),
		c.provideSpecValidator(),
		projectDir,
		moduleName,
	)
}

func (c *Container) provideSpecValidator() *specvalidator.Validator {
	return specvalidator.NewValidator()
}

func (c *Container) provideYamlSpecProvider(projectDir, archFilePath string) *spec.Provider {
	return spec.NewProvider(
		c.provideSourceCode(archFilePath),
		c.provideYamlValidator(projectDir, archFilePath),
	)
}

func (c *Container) provideYamlValidator(projectDir, archFilePath string) *validator.Validator {
	return validator.NewValidator(
		c.provideSourceCodeReferenceResolver(archFilePath),
		c.providePathResolver(),
		projectDir,
	)
}

func (c *Container) providePathResolver() *path.Resolver {
	return path.NewResolver()
}

func (c *Container) provideSourceCode(filePath string) []byte {
	sourceCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to provide source code of archfile: %s", err))
	}

	return sourceCode
}

func (c *Container) provideSourceCodeReferenceResolver(filePath string) *reference.Resolver {
	return reference.NewResolver(
		c.provideSourceCode(filePath),
		filePath,
	)
}

func (c *Container) provideReferenceRender() *code.Render {
	return code.NewRender(
		c.provideColorPrinter(),
	)
}

func (c *Container) provideSpecChecker(projectDir, moduleName string) *checker.Checker {
	return checker.NewChecker(
		c.provideProjectFilesResolver(),
		projectDir,
		moduleName,
	)
}

func (c *Container) provideProjectFilesResolver() *project.Resolver {
	return project.NewResolver()
}
