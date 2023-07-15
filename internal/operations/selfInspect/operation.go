package selfInspect

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Operation struct {
	specAssembler        specAssembler
	projectInfoAssembler projectInfoAssembler
	version              string
}

func NewOperation(
	specAssembler specAssembler,
	projectInfoAssembler projectInfoAssembler,
	version string,
) *Operation {
	return &Operation{
		specAssembler:        specAssembler,
		projectInfoAssembler: projectInfoAssembler,
		version:              version,
	}
}

func (s *Operation) Behave(in models.CmdSelfInspectIn) (models.CmdSelfInspectOut, error) {
	projectInfo, err := s.projectInfoAssembler.ProjectInfo(
		in.ProjectPath,
		in.ArchFile,
	)
	if err != nil {
		return models.CmdSelfInspectOut{}, fmt.Errorf("failed to assemble project info: %w", err)
	}

	spec, err := s.specAssembler.Assemble(projectInfo)
	if err != nil {
		return models.CmdSelfInspectOut{}, fmt.Errorf("failed assemble spec: %w", err)
	}

	return models.CmdSelfInspectOut{
		ModuleName:    projectInfo.ModuleName,
		RootDirectory: projectInfo.Directory,
		LinterVersion: s.version,
		Notices:       s.extractNotices(&spec),
		Suggestions:   s.extractSuggestions(&spec),
	}, nil
}

func (s *Operation) extractNotices(spec *speca.Spec) []models.CmdSelfInspectOutAnnotation {
	return s.asAnnotations(spec.Integrity.DocumentNotices)
}

func (s *Operation) extractSuggestions(spec *speca.Spec) []models.CmdSelfInspectOutAnnotation {
	return s.asAnnotations(spec.Integrity.Suggestions)
}

func (s *Operation) asAnnotations(list []speca.Notice) []models.CmdSelfInspectOutAnnotation {
	annotations := make([]models.CmdSelfInspectOutAnnotation, 0, len(list))

	for _, notice := range list {
		annotations = append(annotations, s.asAnnotation(notice))
	}

	return annotations
}

func (s *Operation) asAnnotation(notice speca.Notice) models.CmdSelfInspectOutAnnotation {
	return models.CmdSelfInspectOutAnnotation{
		Text:      notice.Notice.Error(),
		Reference: notice.Ref,
	}
}
