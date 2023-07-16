package selfInspect

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
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

func (o *Operation) Behave(in models.CmdSelfInspectIn) (models.CmdSelfInspectOut, error) {
	projectInfo, err := o.projectInfoAssembler.ProjectInfo(
		in.ProjectPath,
		in.ArchFile,
	)
	if err != nil {
		return models.CmdSelfInspectOut{}, fmt.Errorf("failed to assemble project info: %w", err)
	}

	spec, err := o.specAssembler.Assemble(projectInfo)
	if err != nil {
		return models.CmdSelfInspectOut{}, fmt.Errorf("failed assemble spec: %w", err)
	}

	return models.CmdSelfInspectOut{
		ModuleName:    projectInfo.ModuleName,
		RootDirectory: projectInfo.Directory,
		LinterVersion: o.version,
		Notices:       o.extractNotices(&spec),
		Suggestions:   o.extractSuggestions(&spec),
	}, nil
}

func (o *Operation) extractNotices(spec *arch.Spec) []models.CmdSelfInspectOutAnnotation {
	return o.asAnnotations(spec.Integrity.DocumentNotices)
}

func (o *Operation) extractSuggestions(spec *arch.Spec) []models.CmdSelfInspectOutAnnotation {
	return o.asAnnotations(spec.Integrity.Suggestions)
}

func (o *Operation) asAnnotations(list []arch.Notice) []models.CmdSelfInspectOutAnnotation {
	annotations := make([]models.CmdSelfInspectOutAnnotation, 0, len(list))

	for _, notice := range list {
		annotations = append(annotations, o.asAnnotation(notice))
	}

	return annotations
}

func (o *Operation) asAnnotation(notice arch.Notice) models.CmdSelfInspectOutAnnotation {
	return models.CmdSelfInspectOutAnnotation{
		Text:      notice.Notice.Error(),
		Reference: notice.Ref,
	}
}
