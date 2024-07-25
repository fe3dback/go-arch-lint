package container

import (
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
			c.serviceConfigValidatorCmnComponents(),
		)
	})
}

func (c *Container) serviceConfigValidatorCmnComponents() *validator.CommonComponentsValidator {
	return once(validator.NewCommonComponentsValidator)
}
