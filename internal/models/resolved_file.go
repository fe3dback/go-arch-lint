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

	ResolvedFile struct {
		Path    string
		Imports []ResolvedImport
	}
)
