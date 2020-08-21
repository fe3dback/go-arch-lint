package checker

import (
	"fmt"
	"io"
	"strings"

	"github.com/fe3dback/go-arch-lint/files"
	"github.com/fe3dback/go-arch-lint/spec"
)

type Checker struct {
	rootDirectory string
	arch          spec.Arch
	projectFiles  files.ResolveResult
	errorsLog     io.Writer
}

func NewChecker(
	rootDirectory string,
	arch *spec.Arch,
	projectFiles files.ResolveResult,
	errorsLog io.Writer,
) *Checker {
	return &Checker{
		rootDirectory: rootDirectory,
		arch:          *arch,
		projectFiles:  projectFiles,
		errorsLog:     errorsLog,
	}
}

func (arc *Checker) Check() {
	for _, file := range arc.projectFiles.Files {
		component := arc.resolveComponent(file.Path)
		if component == nil {
			arc.logWarning(fmt.Sprintf("File '%s' not attached to any component in archfile", file.Path))
			continue
		}

		arc.checkFile(component, file)
	}
}

func (arc *Checker) resolveComponent(filePath string) *spec.Component {
	matched := make(map[string]*spec.Component)

	for _, component := range arc.arch.Components {
		for _, componentDirectory := range component.ResolvedPaths {
			if strings.HasPrefix(filePath, componentDirectory.AbsPath) {
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
		allowed := arc.checkImport(component, resolvedImport)
		if !allowed {
			arc.logWarning(fmt.Sprintf("Component '%s' file '%s' shouldnot depend on '%s'",
				component.Name,
				strings.TrimPrefix(file.Path, arc.rootDirectory),
				resolvedImport.Name,
			))
			continue
		}
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
	// todo
	if arc.arch.Allow.DepOnAnyVendor {
		return true
	}

	if component.SpecialFlags.AllowAllVendorDeps {
		return true
	}

	return false
}

func (arc *Checker) checkProjectImport(component *spec.Component, resolvedImport files.ResolvedImport) bool {
	if component.SpecialFlags.AllowAllProjectDeps {
		return true
	}

	for _, allowedImport := range component.AllowedImports {
		if allowedImport.ImportPath == resolvedImport.Name {
			return true
		}
	}

	return false
}

func (arc *Checker) logWarning(warn string) {
	msg := fmt.Sprintf("[WARNING] %s\n", warn)
	_, _ = arc.errorsLog.Write([]byte(msg))
}
