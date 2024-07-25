package config

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	Reader interface {
		Read(path models.PathAbsolute) (models.Config, error)
	}

	Validator interface {
		Validate(config models.Config) error
	}

	Assembler interface {
		// todo: add spec assembler
	}

	// todo: maybe another business validator here?
)
