package check

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Service struct {
	specAssembler   SpecAssembler
	referenceRender ReferenceRender
}

func NewService(
	specAssembler SpecAssembler,
	referenceRender ReferenceRender,
) *Service {
	return &Service{
		specAssembler:   specAssembler,
		referenceRender: referenceRender,
	}
}

func (s *Service) Behave() (models.Check, error) {
	spec, err := s.specAssembler.Assemble()
	if err != nil {
		return models.Check{}, fmt.Errorf("failed to assemble spec: %w", err)
	}

	model := models.Check{
		ModuleName:             spec.ModuleName.Value(),
		DocumentNotices:        s.assembleNotice(spec.Integrity),
		ArchHasWarnings:        false,
		ArchWarningsDependency: []models.CheckArchWarningDependency{},
		ArchWarningsMatch:      []models.CheckArchWarningMatch{},
	}

	return model, nil
}

func (s *Service) assembleNotice(integrity speca.Integrity) []models.CheckNotice {
	notices := make([]speca.Notice, 0)
	notices = append(notices, integrity.DocumentNotices...)
	notices = append(notices, integrity.SpecNotices...)

	results := make([]models.CheckNotice, 0)
	for _, notice := range notices {
		results = append(results, models.CheckNotice{
			Text:              fmt.Sprintf("%s", notice.Notice),
			File:              notice.Ref.File,
			Line:              notice.Ref.Line,
			Offset:            notice.Ref.Offset,
			SourceCodePreview: s.referenceRender.SourceCode(notice.Ref, 1, true),
		})
	}

	return results
}
