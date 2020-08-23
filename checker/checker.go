package checker

import (
	"fmt"
	"path/filepath"
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

	return arc.longestPathComponent(matched)
}

func (arc *Checker) longestPathComponent(matched map[string]*spec.Component) *spec.Component {
	longest := ""
	var targetComponent *spec.Component

	for path, component := range matched {
		if len(path) > len(longest) {
			longest = path
			targetComponent = component
		}
	}

	return targetComponent
}

func (arc *Checker) checkFile(component *spec.Component, file *files.ResolvedFile) {
	for _, resolvedImport := range file.Imports {
		if arc.checkImport(component, resolvedImport) {
			continue
		}

		arc.result.addDependencyWarning(WarningDep{
			ComponentName:      component.Name,
			FileRelativePath:   strings.TrimPrefix(file.Path, arc.rootDirectory),
			FileAbsolutePath:   file.Path,
			ResolvedImportName: resolvedImport.Name,
		})
		break
	}
}

func (arc *Checker) checkImport(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	switch resolvedImport.ImportType {
	case files.ImportTypeStdLib:
		return true
	case files.ImportTypeVendor:
		return arc.checkVendorImport(component, resolvedImport)
	case files.ImportTypeProject:
		return arc.checkProjectImport(component, resolvedImport)
	default:
		panic(fmt.Sprintf("unknown import type: %+v", resolvedImport))
	}
}

func (arc *Checker) checkVendorImport(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	if arc.arch.Allow.DepOnAnyVendor {
		return true
	}

	if component.SpecialFlags.AllowAllVendorDeps {
		return true
	}

	return arc.checkImportPath(component, resolvedImport)
}

func (arc *Checker) checkProjectImport(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllProjectDeps {
		return true
	}

	return arc.checkImportPath(component, resolvedImport)
}

func (arc *Checker) checkImportPath(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	for _, allowedImport := range component.AllowedImports {
		if allowedImport.ImportPath == resolvedImport.Name {
			return true
		}
	}

	return false
}
