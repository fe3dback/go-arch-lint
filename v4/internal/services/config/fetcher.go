package config

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Fetcher struct {
	Reader    Reader
	Validator Validator
	Assembler Assembler
}

func NewFetcher(reader Reader, validator Validator, assembler Assembler) *Fetcher {
	return &Fetcher{Reader: reader, Validator: validator, Assembler: assembler}
}

func (f *Fetcher) FetchSpec(path models.PathAbsolute) (models.Config, error) {
	// todo: change return type to models.Spec, error

	conf, err := f.Reader.Read(path)
	if err != nil {
		return models.Config{}, fmt.Errorf("failed read config at '%s': %w", path, err)
	}

	err = f.Validator.Validate(conf)
	if err != nil {
		return models.Config{}, fmt.Errorf("invalid config: %w", err)
	}

	// todo: spec assemble

	return conf, nil
}
