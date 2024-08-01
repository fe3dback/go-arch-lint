package config

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Fetcher struct {
	Reader    reader
	Validator validator
	Assembler assembler
}

func NewFetcher(reader reader, validator validator, assembler assembler) *Fetcher {
	return &Fetcher{Reader: reader, Validator: validator, Assembler: assembler}
}

func (f *Fetcher) FetchSpec(path models.PathAbsolute) (models.Spec, error) {
	conf, err := f.Reader.Read(path)
	if err != nil {
		return models.Spec{}, fmt.Errorf("failed read config at '%s': %w", path, err)
	}

	err = f.Validator.Validate(conf)
	if err != nil {
		return models.Spec{}, fmt.Errorf("invalid config: %w", err)
	}

	spec, err := f.Assembler.Assemble(conf)
	if err != nil {
		return models.Spec{}, fmt.Errorf("failed assemble spec: %w", err)
	}

	return spec, nil
}
