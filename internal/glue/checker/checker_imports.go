package checker

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Imports struct {
	spec                 speca.Spec
	projectFilesResolver ProjectFilesResolver
	result               results
}

func NewImport(
	projectFilesResolver ProjectFilesResolver,
) *Imports {
	return &Imports{
		result:               newResults(),
		projectFilesResolver: projectFilesResolver,
	}
}

func (c *Imports) Check(spec speca.Spec) (models.CheckResult, error) {
	c.spec = spec

	projectFiles, err := c.projectFilesResolver.ProjectFiles(spec)
	if err != nil {
		return models.CheckResult{}, fmt.Errorf("failed to resolve project files: %w", err)
	}

	components := c.assembleComponentsMap(spec)

	for _, projectFile := range projectFiles {
		if projectFile.ComponentID == nil {
			c.result.addNotMatchedWarning(models.CheckArchWarningMatch{
				Reference:        speca.NewEmptyReference(),
				FileRelativePath: strings.TrimPrefix(projectFile.File.Path, spec.RootDirectory.Value()),
				FileAbsolutePath: projectFile.File.Path,
			})

			continue
		}

		componentID := *projectFile.ComponentID
		if component, ok := components[componentID]; ok {
			err := c.checkFile(component, projectFile.File)
			if err != nil {
				return models.CheckResult{}, fmt.Errorf("failed check file '%s': %w", projectFile.File.Path, err)
			}

			continue
		}

		return models.CheckResult{}, fmt.Errorf("not found component '%s' in map", componentID)
	}

	return c.result.assembleSortedResults(), nil
}

func (c *Imports) assembleComponentsMap(spec speca.Spec) map[string]speca.Component {
	results := make(map[string]speca.Component)

	for _, component := range spec.Components {
		results[component.Name.Value()] = component
	}

	return results
}

func (c *Imports) checkFile(component speca.Component, file models.ProjectFile) error {
	for _, resolvedImport := range file.Imports {
		allowed, err := checkImport(component, resolvedImport, c.spec.Allow.DepOnAnyVendor.Value())
		if err != nil {
			return fmt.Errorf("failed check import '%s': %w",
				resolvedImport.Name,
				err,
			)
		}

		if allowed {
			continue
		}

		c.result.addDependencyWarning(models.CheckArchWarningDependency{
			Reference:          component.Name.Reference(),
			ComponentName:      component.Name.Value(),
			FileRelativePath:   strings.TrimPrefix(file.Path, c.spec.RootDirectory.Value()),
			FileAbsolutePath:   file.Path,
			ResolvedImportName: resolvedImport.Name,
		})
	}

	return nil
}

func checkImport(
	component speca.Component,
	resolvedImport models.ResolvedImport,
	allowDependOnAnyVendor bool,
) (bool, error) {
	switch resolvedImport.ImportType {
	case models.ImportTypeStdLib:
		return true, nil
	case models.ImportTypeVendor:
		if allowDependOnAnyVendor {
			return true, nil
		}

		return checkVendorImport(component, resolvedImport)
	case models.ImportTypeProject:
		return checkProjectImport(component, resolvedImport), nil
	default:
		panic(fmt.Sprintf("unknown import type: %+v", resolvedImport))
	}
}

func checkVendorImport(component speca.Component, resolvedImport models.ResolvedImport) (bool, error) {
	if component.SpecialFlags.AllowAllVendorDeps.Value() {
		return true, nil
	}

	for _, vendorGlob := range component.AllowedVendorGlobs {
		matched, err := vendorGlob.Value().Match(resolvedImport.Name)
		if err != nil {
			return false, models.NewReferableErr(
				fmt.Errorf("invalid vendor glob '%s': %w",
					string(vendorGlob.Value()),
					err,
				),
				vendorGlob.Reference(),
			)
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func checkProjectImport(component speca.Component, resolvedImport models.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllProjectDeps.Value() {
		return true
	}

	for _, allowedImportRef := range component.AllowedProjectImports {
		allowedImport := allowedImportRef.Value()

		if allowedImport.ImportPath == resolvedImport.Name {
			return true
		}
	}

	return false
}
