package models

import (
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

const (
	ImportTypeStdLib ImportType = iota
	ImportTypeProject
	ImportTypeVendor
)

type (
	ImportType uint8

	FileHold struct {
		File        ProjectFile
		ComponentID *string
	}

	ProjectFile struct {
		Path    string
		Imports []ResolvedImport
	}

	ResolvedImport struct {
		Name       string
		ImportType ImportType
		Reference  common.Reference
	}
)

// GetImportType classifies an import path as std, project, or vendor
func GetImportType(importPath string, moduleName string, stdPackages map[string]struct{}) ImportType {
	if _, ok := stdPackages[importPath]; ok {
		return ImportTypeStdLib
	}

	// We can't use a straight prefix match here because the module name could be a substring of the import path.
	// For example, if the module name is "example.com/foo/bar", we do not want to match "example.com/foo/bar-utils"
	if importPath == moduleName || strings.HasPrefix(importPath, moduleName+"/") {
		return ImportTypeProject
	}

	return ImportTypeVendor
}
