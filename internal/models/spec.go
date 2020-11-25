package models

import (
	"regexp"
)

type (
	ArchSpec struct {
		RootDirectory       string
		ModuleName          string
		Allow               ArchAllow
		Components          []*Component
		Exclude             []*ResolvedPath
		ExcludeFilesMatcher []*regexp.Regexp
	}

	ArchAllow struct {
		DepOnAnyVendor bool
	}

	ComponentName = string
	VendorName    = string
	Component     struct {
		Name           ComponentName
		LocalPathMask  string
		ResolvedPaths  []*ResolvedPath
		MayDependOn    []ComponentName
		CanUse         []VendorName
		AllowedImports []*ResolvedPath
		SpecialFlags   *SpecialFlags
	}

	SpecialFlags struct {
		AllowAllProjectDeps bool
		AllowAllVendorDeps  bool
	}
)
