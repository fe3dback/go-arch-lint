package check

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Service struct {
	specAssembler SpecAssembler
}

func NewService(
	specAssembler SpecAssembler,
) *Service {
	return &Service{
		specAssembler: specAssembler,
	}
}

func (s *Service) Behave() (models.Check, error) {
	spec, err := s.specAssembler.Assemble()
	if err != nil {
		return models.Check{}, fmt.Errorf("failed to assemble spec: %w", err)
	}

	_ = spec // todo

	return models.Check{}, nil
}
