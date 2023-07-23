package models

import "github.com/fe3dback/go-arch-lint/internal/models/common"

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
