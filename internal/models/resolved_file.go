package models

const (
	ImportTypeStdLib ImportType = iota
	ImportTypeProject
	ImportTypeVendor
)

type (
	ImportType uint8

	ResolvedImport struct {
		Name       string
		ImportType ImportType
	}

	ProjectFile struct {
		Path    string
		Imports []ResolvedImport
	}

	FileHold struct {
		File        ProjectFile
		ComponentID *string
	}
)
