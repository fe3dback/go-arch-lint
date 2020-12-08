package checker

import (
	"fmt"
	"path/filepath"
	"sort"
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

	projectFiles, err := c.projectFilesResolver.Resolve(
		c.rootDirectory,
		c.moduleName,
		refPathToList(spec.Exclude),
		refRegExpToList(spec.ExcludeFilesMatcher),
	)
	if err != nil {
		return models.CheckResult{}, fmt.Errorf("failed to resolve project files: %w", err)
	}

	for _, file := range projectFiles {
		component := c.resolveComponent(file.Path)
		if component == nil {
			c.result.addNotMatchedWarning(models.CheckArchWarningMatch{
				Reference:        speca.NewEmptyReference(),
				FileRelativePath: strings.TrimPrefix(file.Path, c.rootDirectory),
				FileAbsolutePath: file.Path,
			})

			continue
		}

		c.checkFile(component, file)
	}

	return c.result.assembleSortedResults(), nil
}

func (c *Checker) resolveComponent(filePath string) *speca.Component {
	matched := make(map[string]*speca.Component)
	directory := filepath.Dir(filePath)

	for _, component := range c.spec.Components {
		component := component

		for _, componentDirectoryRef := range component.ResolvedPaths {
			componentDirectory := componentDirectoryRef.Value()

			if strings.HasPrefix(directory, componentDirectory.AbsPath) {
				suffixPath := strings.TrimPrefix(directory, componentDirectory.AbsPath)

				if strings.Contains(suffixPath, "/") {
					continue
				}

				matched[componentDirectory.ImportPath] = &component
				continue
			}
		}
	}

	return longestPathComponent(matched)
}

func (c *Checker) checkFile(component *speca.Component, file models.ResolvedFile) {
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

func longestPathComponent(matched map[string]*speca.Component) *speca.Component {
	// work only with keys
	sortedPaths := make([]string, len(matched))
	for path := range matched {
		sortedPaths = append(sortedPaths, path)
	}

	sort.Strings(sortedPaths)

	// find longest
	longest := ""
	for _, path := range sortedPaths {
		if len(path) > len(longest) {
			longest = path
		}
	}

	if longest == "" {
		return nil
	}

	return matched[longest]
}

func checkImport(
	component *speca.Component,
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

func checkVendorImport(component *speca.Component, resolvedImport models.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllVendorDeps.Value() {
		return true
	}

	return checkImportPath(component, resolvedImport)
}

func checkProjectImport(component *speca.Component, resolvedImport models.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllProjectDeps.Value() {
		return true
	}

	return checkImportPath(component, resolvedImport)
}

func checkImportPath(component *speca.Component, resolvedImport models.ResolvedImport) bool {
	for _, allowedImportRef := range component.AllowedImports {
		allowedImport := allowedImportRef.Value()

		if allowedImport.ImportPath == resolvedImport.Name {
			return true
		}
	}

	return false
}
