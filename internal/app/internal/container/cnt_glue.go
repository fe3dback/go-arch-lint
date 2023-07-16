package container

import (
	"github.com/fe3dback/go-arch-lint/internal/services/checker"
	"github.com/fe3dback/go-arch-lint/internal/services/code"
	"github.com/fe3dback/go-arch-lint/internal/services/path"
	"github.com/fe3dback/go-arch-lint/internal/services/project/holder"
	"github.com/fe3dback/go-arch-lint/internal/services/project/info"
	"github.com/fe3dback/go-arch-lint/internal/services/project/resolver"
	"github.com/fe3dback/go-arch-lint/internal/services/project/scanner"
	"github.com/fe3dback/go-arch-lint/internal/services/schema"
	specassembler "github.com/fe3dback/go-arch-lint/internal/services/spec/assembler"
	specvalidator "github.com/fe3dback/go-arch-lint/internal/services/spec/validator"
	"github.com/fe3dback/go-arch-lint/internal/services/yaml/reference"
	"github.com/fe3dback/go-arch-lint/internal/services/yaml/spec"
)

func (c *Container) provideSpecAssembler() *specassembler.Assembler {
	return specassembler.NewAssembler(
		c.provideYamlSpecProvider(),
		c.provideSpecValidator(),
		c.providePathResolver(),
	)
}

func (c *Container) provideSpecValidator() *specvalidator.Validator {
	return specvalidator.NewValidator(
		c.providePathResolver(),
	)
}

func (c *Container) provideYamlSpecProvider() *spec.Provider {
	return spec.NewProvider(
		c.provideSourceCodeReferenceResolver(),
		c.provideJsonSchemaProvider(),
	)
}

func (c *Container) providePathResolver() *path.Resolver {
	return path.NewResolver()
}

func (c *Container) provideSourceCodeReferenceResolver() *reference.Resolver {
	return reference.NewResolver()
}

func (c *Container) provideReferenceRender() *code.Render {
	return code.NewRender(
		c.provideColorPrinter(),
	)
}

func (c *Container) provideSpecChecker() *checker.CompositeChecker {
	return checker.NewCompositeChecker(
		c.provideSpecImportsChecker(),
		c.provideSpecDeepScanChecker(),
	)
}

func (c *Container) provideSpecImportsChecker() *checker.Imports {
	return checker.NewImport(
		c.provideProjectFilesResolver(),
	)
}

func (c *Container) provideSpecDeepScanChecker() *checker.DeepScan {
	return checker.NewDeepScan(
		c.provideProjectFilesResolver(),
		c.provideReferenceRender(),
	)
}

func (c *Container) provideProjectFilesResolver() *resolver.Resolver {
	return resolver.NewResolver(
		c.provideProjectFilesScanner(),
		c.provideProjectFilesHolder(),
	)
}

func (c *Container) provideProjectFilesScanner() *scanner.Scanner {
	return scanner.NewScanner()
}

func (c *Container) provideProjectFilesHolder() *holder.Holder {
	return holder.NewHolder()
}

func (c *Container) provideProjectInfoAssembler() *info.Assembler {
	return info.NewAssembler()
}

func (c *Container) provideJsonSchemaProvider() *schema.Provider {
	return schema.NewProvider()
}
