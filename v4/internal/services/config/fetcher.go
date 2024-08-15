package config

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Fetcher struct {
	reader        reader
	validator     validator
	assembler     assembler
	moduleFetcher moduleFetcher
}

func NewFetcher(reader reader, validator validator, assembler assembler, moduleFetcher moduleFetcher) *Fetcher {
	return &Fetcher{reader: reader, validator: validator, assembler: assembler, moduleFetcher: moduleFetcher}
}

func (f *Fetcher) FetchSpec() (models.Spec, error) {
	project, err := f.moduleFetcher.Fetch()
	if err != nil {
		return models.Spec{}, fmt.Errorf("failed fetch module info: %w", err)
	}

	conf, err := f.reader.Read(project.ConfigPath)
	if err != nil {
		return models.Spec{}, fmt.Errorf("failed read config at '%s': %w", project.ConfigPath, err)
	}

	err = f.validator.Validate(conf)
	if err != nil {
		return models.Spec{}, fmt.Errorf("invalid config: %w", err)
	}

	spec, err := f.assembler.Assemble(conf)
	if err != nil {
		return models.Spec{}, fmt.Errorf("failed assemble spec: %w", err)
	}

	return spec, nil
}
