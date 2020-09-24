package checker

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fe3dback/go-arch-lint/files"
	"github.com/fe3dback/go-arch-lint/spec"
)

type Checker struct {
	rootDirectory string
	arch          spec.Arch
	projectFiles  files.ResolveResult
	result        *CheckResult
}

func NewChecker(
	rootDirectory string,
	arch *spec.Arch,
	projectFiles files.ResolveResult,
) *Checker {
	return &Checker{
		rootDirectory: rootDirectory,
		arch:          *arch,
		projectFiles:  projectFiles,
		result:        newCheckResult(),
	}
}

func (arc *Checker) Check() CheckResult {
	for _, file := range arc.projectFiles.Files {
		component := arc.resolveComponent(file.Path)
		if component == nil {
			arc.result.addNotMatchedWarning(WarningNotMatched{
				FileRelativePath: strings.TrimPrefix(file.Path, arc.rootDirectory),
				FileAbsolutePath: file.Path,
			})

			continue
		}

		arc.checkFile(component, file)
	}

	return *arc.result
}

func (arc *Checker) resolveComponent(filePath string) *spec.Component {
	matched := make(map[string]*spec.Component)
	directory := filepath.Dir(filePath)

	for _, component := range arc.arch.Components {
		for _, componentDirectory := range component.ResolvedPaths {
			if strings.HasPrefix(directory, componentDirectory.AbsPath) {
				suffixPath := strings.TrimPrefix(directory, componentDirectory.AbsPath)

				if strings.Contains(suffixPath, "/") {
					continue
				}

				matched[componentDirectory.ImportPath] = component
				continue
			}
		}
	}

	return longestPathComponent(matched)
}

func (arc *Checker) checkFile(component *spec.Component, file *files.ResolvedFile) {
	for _, resolvedImport := range file.Imports {
		if checkImport(component, resolvedImport, arc.arch.Allow.DepOnAnyVendor) {
			continue
		}

		arc.result.addDependencyWarning(WarningDep{
			ComponentName:      component.Name,
			FileRelativePath:   strings.TrimPrefix(file.Path, arc.rootDirectory),
			FileAbsolutePath:   file.Path,
			ResolvedImportName: resolvedImport.Name,
		})
	}
}

func longestPathComponent(matched map[string]*spec.Component) *spec.Component {
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
	component *spec.Component,
	resolvedImport files.ResolvedImport,
	allowDependOnAnyVendor bool,
) bool {
	switch resolvedImport.ImportType {
	case files.ImportTypeStdLib:
		return true
	case files.ImportTypeVendor:
		if allowDependOnAnyVendor {
			return true
		}

		return checkVendorImport(component, resolvedImport)
	case files.ImportTypeProject:
		return checkProjectImport(component, resolvedImport)
	default:
		panic(fmt.Sprintf("unknown import type: %+v", resolvedImport))
	}
}

func checkVendorImport(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllVendorDeps {
		return true
	}

	return checkImportPath(component, resolvedImport)
}

func checkProjectImport(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllProjectDeps {
		return true
	}

	return checkImportPath(component, resolvedImport)
}

func checkImportPath(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	for _, allowedImport := range component.AllowedImports {
		if allowedImport.ImportPath == resolvedImport.Name {
			return true
		}
	}

	return false
}
