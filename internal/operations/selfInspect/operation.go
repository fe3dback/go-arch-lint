package selfInspect

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Operation struct {
	specAssembler specAssembler
	version       string
}

func NewOperation(specAssembler specAssembler, version string) *Operation {
	return &Operation{
		specAssembler: specAssembler,
		version:       version,
	}
}

func (s *Operation) Behave() (models.SelfInspect, error) {
	spec, err := s.specAssembler.Assemble()
	if err != nil {
		return models.SelfInspect{}, fmt.Errorf("failed assemble spec: %w", err)
	}

	return models.SelfInspect{
		LinterVersion: s.version,
		Notices:       s.extractNotices(&spec),
		Suggestions:   s.extractSuggestions(&spec),
	}, nil
}

func (s *Operation) extractNotices(spec *speca.Spec) []models.Annotation {
	return s.asAnnotations(spec.Integrity.DocumentNotices)
}

func (s *Operation) extractSuggestions(spec *speca.Spec) []models.Annotation {
	return s.asAnnotations(spec.Integrity.Suggestions)
}

func (s *Operation) asAnnotations(list []speca.Notice) []models.Annotation {
	annotations := make([]models.Annotation, 0, len(list))

	for _, notice := range list {
		annotations = append(annotations, s.asAnnotation(notice))
	}

	return annotations
}

func (s *Operation) asAnnotation(notice speca.Notice) models.Annotation {
	return models.Annotation{
		Text:      notice.Notice.Error(),
		Reference: notice.Ref,
	}
}
