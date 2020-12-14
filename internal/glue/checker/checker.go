package checker

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Checker struct {
	spec                 speca.Spec
	projectFilesResolver ProjectFilesResolver
	result               results
	rootDirectory        string
	moduleName           string
}

func NewChecker(
	projectFilesResolver ProjectFilesResolver,
	rootDirectory string,
	moduleName string,
) *Checker {
	return &Checker{
		result:               newResults(),
		projectFilesResolver: projectFilesResolver,
		rootDirectory:        rootDirectory,
		moduleName:           moduleName,
	}
}

func (c *Checker) Check(spec speca.Spec) (models.CheckResult, error) {
	c.spec = spec

	projectFiles, err := c.projectFilesResolver.ProjectFiles(
		c.rootDirectory,
		c.moduleName,
		spec,
	)
	if err != nil {
		return models.CheckResult{}, fmt.Errorf("failed to resolve project files: %w", err)
	}

	components := c.assembleComponentsMap(spec)

	for _, projectFile := range projectFiles {
		if projectFile.ComponentID == nil {
			c.result.addNotMatchedWarning(models.CheckArchWarningMatch{
				Reference:        speca.NewEmptyReference(),
				FileRelativePath: strings.TrimPrefix(projectFile.File.Path, c.rootDirectory),
				FileAbsolutePath: projectFile.File.Path,
			})

			continue
		}

		componentID := *projectFile.ComponentID
		if component, ok := components[componentID]; ok {
			c.checkFile(component, projectFile.File)

			continue
		}

		return models.CheckResult{}, fmt.Errorf("not found component '%s' in map", componentID)
	}

	return c.result.assembleSortedResults(), nil
}

func (c *Checker) assembleComponentsMap(spec speca.Spec) map[string]speca.Component {
	results := make(map[string]speca.Component)

	for _, component := range spec.Components {
		results[component.Name.Value()] = component
	}

	return results
}

func (c *Checker) checkFile(component speca.Component, file models.ProjectFile) {
	for _, resolvedImport := range file.Imports {
		if checkImport(component, resolvedImport, c.spec.Allow.DepOnAnyVendor.Value()) {
			continue
		}

		c.result.addDependencyWarning(models.CheckArchWarningDependency{
			Reference:          component.Name.Reference(),
			ComponentName:      component.Name.Value(),
			FileRelativePath:   strings.TrimPrefix(file.Path, c.rootDirectory),
			FileAbsolutePath:   file.Path,
			ResolvedImportName: resolvedImport.Name,
		})
	}
}

func checkImport(
	component speca.Component,
	resolvedImport models.ResolvedImport,
	allowDependOnAnyVendor bool,
) bool {
	switch resolvedImport.ImportType {
	case models.ImportTypeStdLib:
		return true
	case models.ImportTypeVendor:
		if allowDependOnAnyVendor {
			return true
		}

		return checkVendorImport(component, resolvedImport)
	case models.ImportTypeProject:
		return checkProjectImport(component, resolvedImport)
	default:
		panic(fmt.Sprintf("unknown import type: %+v", resolvedImport))
	}
}

func checkVendorImport(component speca.Component, resolvedImport models.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllVendorDeps.Value() {
		return true
	}

	return checkImportPath(component, resolvedImport)
}

func checkProjectImport(component speca.Component, resolvedImport models.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllProjectDeps.Value() {
		return true
	}

	return checkImportPath(component, resolvedImport)
}

func checkImportPath(component speca.Component, resolvedImport models.ResolvedImport) bool {
	for _, allowedImportRef := range component.AllowedImports {
		allowedImport := allowedImportRef.Value()

		if allowedImport.ImportPath == resolvedImport.Name {
			return true
		}
	}

	return false
}
