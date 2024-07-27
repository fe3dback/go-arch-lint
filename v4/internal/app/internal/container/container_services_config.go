package container

import (
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/config"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/reader"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/reader/yaml"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/config/validator"
)

func (c *Container) serviceConfigFetcher() *config.Fetcher {
	return once(func() *config.Fetcher {
		return config.NewFetcher(
			c.serviceConfigReader(),
			c.serviceConfigValidator(),
			nil, // todo: assembler
		)
	})
}

func (c *Container) serviceConfigReader() *reader.Reader {
	return once(func() *reader.Reader {
		return reader.NewReader(
			c.serviceConfigReaderYAML(),
		)
	})
}

func (c *Container) serviceConfigReaderYAML() *yaml.Reader {
	return once(func() *yaml.Reader {
		return yaml.NewReader()
	})
}

func (c *Container) serviceConfigValidator() *validator.Root {
	return once(func() *validator.Root {
		return validator.NewRoot(
			c.cCtx.Bool(models.FlagSkipMissUsages),
			c.serviceConfigValidatorWorkdir(),
			c.serviceConfigValidatorExcludedFiles(),
			c.serviceConfigValidatorCmnComponents(),
			c.serviceConfigValidatorCmnVendors(),
			c.serviceConfigValidatorComponents(),
			c.serviceConfigValidatorDeps(),
			c.serviceConfigValidatorDepsComponents(),
			c.serviceConfigValidatorDepsVendors(),
		)
	})
}

func (c *Container) serviceConfigValidatorWorkdir() *validator.WorkdirValidator {
	return once(func() *validator.WorkdirValidator {
		return validator.NewWorkdirValidator(
			c.serviceProjectPathHelper(),
		)
	})
}

func (c *Container) serviceConfigValidatorExcludedFiles() *validator.ExcludedFilesValidator {
	return once(validator.NewExcludedFilesValidator)
}

func (c *Container) serviceConfigValidatorCmnComponents() *validator.CommonComponentsValidator {
	return once(validator.NewCommonComponentsValidator)
}

func (c *Container) serviceConfigValidatorCmnVendors() *validator.CommonVendorsValidator {
	return once(validator.NewCommonVendorsValidator)
}

func (c *Container) serviceConfigValidatorComponents() *validator.ComponentsValidator {
	return once(func() *validator.ComponentsValidator {
		return validator.NewComponentsValidator(
			c.serviceProjectPathHelper(),
		)
	})
}

func (c *Container) serviceConfigValidatorDeps() *validator.DepsValidator {
	return once(func() *validator.DepsValidator {
		return validator.NewDepsValidator(
			c.serviceProjectPathHelper(),
		)
	})
}

func (c *Container) serviceConfigValidatorDepsComponents() *validator.DepsComponentsValidator {
	return once(func() *validator.DepsComponentsValidator {
		return validator.NewDepsComponentsValidator(
			c.serviceProjectPathHelper(),
		)
	})
}

func (c *Container) serviceConfigValidatorDepsVendors() *validator.DepsVendorsValidator {
	return once(func() *validator.DepsVendorsValidator {
		return validator.NewDepsVendorsValidator(
			c.serviceProjectPathHelper(),
		)
	})
}
