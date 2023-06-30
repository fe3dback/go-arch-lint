package models

type (
	FlagsGraph struct {
		Project        ProjectInfo
		OutFile        string
		Type           GraphType
		IncludeVendors bool
		Focus          string
	}
)
