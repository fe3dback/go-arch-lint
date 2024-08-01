package config

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	reader interface {
		Read(path models.PathAbsolute) (models.Config, error)
	}

	validator interface {
		Validate(config models.Config) error
	}

	assembler interface {
		Assemble(conf models.Config) (models.Spec, error)
	}
)
