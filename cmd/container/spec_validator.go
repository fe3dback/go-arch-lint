package container

import (
	"github.com/fe3dback/go-arch-lint/spec/annotated_validator"
	"github.com/fe3dback/go-arch-lint/spec/validator"
	"github.com/fe3dback/go-arch-lint/spec/warnparser"
)

func (c *Container) provideSpecValidator() *validator.ArchFileValidator {
	return validator.NewArchFileValidator(
		c.providePathResolver(),
		c.provideArchSpec(),
		c.provideProjectRootDirectory(),
	)
}

func (c *Container) provideSpecWarnParser() *warnparser.WarningSourceParser {
	return warnparser.NewWarningSourceParser()
}

func (c *Container) ProvideSpecAnnotatedValidator() *annotated_validator.AnnotatedValidator {
	return annotated_validator.NewAnnotatedValidator(
		c.provideSpecValidator(),
		c.provideSpecWarnParser(),
		c.provideArchFileSourceCode(),
		c.provideFlagUseColors(),
	)
}
