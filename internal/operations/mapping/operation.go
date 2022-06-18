package mapping

import (
	"context"
	"fmt"
	"sort"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Operation struct {
	specAssembler        SpecAssembler
	projectFilesResolver ProjectFilesResolver
}

func NewOperation(
	specAssembler SpecAssembler,
	projectFilesResolver ProjectFilesResolver,
) *Operation {
	return &Operation{
		specAssembler:        specAssembler,
		projectFilesResolver: projectFilesResolver,
	}
}

func (s *Operation) Behave(ctx context.Context, scheme models.MappingScheme) (models.Mapping, error) {
	spec, err := s.specAssembler.Assemble()
	if err != nil {
		return models.Mapping{}, fmt.Errorf("failed to assemble spec: %w", err)
	}

	projectFiles, err := s.projectFilesResolver.ProjectFiles(ctx, spec)
	if err != nil {
		return models.Mapping{}, fmt.Errorf("failed to resolve project files: %w", err)
	}

	return models.Mapping{
		ProjectDirectory: spec.RootDirectory.Value(),
		ModuleName:       spec.ModuleName.Value(),
		MappingGrouped:   assembleMappingByComponent(spec, projectFiles),
		MappingList:      assembleMappingByFile(projectFiles),
		Scheme:           scheme,
	}, nil
}

func assembleMappingByComponent(
	spec speca.Spec,
	projectFiles []models.FileHold,
) []models.MappingGrouped {
	tmp := make(map[string]*models.MappingGrouped)

	for _, projectFile := range projectFiles {
		componentName := componentName(projectFile.ComponentID)
		if _, exist := tmp[componentName]; !exist {
			tmp[componentName] = &models.MappingGrouped{
				ComponentName: componentName,
				FileNames:     []string{},
			}
		}

		fileName := projectFile.File.Path
		tmp[componentName].FileNames = append(
			tmp[componentName].FileNames,
			fileName,
		)
	}

	mapping := make([]models.MappingGrouped, 0)
	for _, component := range spec.Components {
		componentName := component.Name.Value()
		if grouped, exist := tmp[componentName]; exist {
			sort.Strings(grouped.FileNames)
			mapping = append(mapping, *grouped)
			continue
		}

		mapping = append(mapping, models.MappingGrouped{
			ComponentName: componentName,
			FileNames:     []string{},
		})
	}

	emptyComponentID := componentName(nil)
	if _, hasNotAttached := tmp[emptyComponentID]; hasNotAttached {
		notAttachedFiles := tmp[emptyComponentID].FileNames

		if len(notAttachedFiles) > 0 {
			sort.Strings(notAttachedFiles)
			mapping = append(mapping, models.MappingGrouped{
				ComponentName: emptyComponentID,
				FileNames:     notAttachedFiles,
			})
		}
	}

	sort.Slice(mapping, func(i, j int) bool {
		return mapping[i].ComponentName < mapping[j].ComponentName
	})

	return mapping
}

func assembleMappingByFile(projectFiles []models.FileHold) []models.MappingList {
	mapping := make([]models.MappingList, 0)
	exist := make(map[string]struct{})

	for _, projectFile := range projectFiles {
		fileName := projectFile.File.Path

		if _, exist := exist[fileName]; exist {
			continue
		}

		mapping = append(mapping, models.MappingList{
			FileName:      fileName,
			ComponentName: componentName(projectFile.ComponentID),
		})

		exist[fileName] = struct{}{}
	}

	sort.Slice(mapping, func(i, j int) bool {
		return mapping[i].FileName < mapping[j].FileName
	})

	return mapping
}

func componentName(id *string) string {
	if id == nil {
		return "[not attached]"
	}

	return *id
}
