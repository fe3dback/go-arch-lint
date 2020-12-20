package speca

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	Spec struct {
		RootDirectory       ReferableString
		ModuleName          ReferableString
		Allow               Allow
		Components          []Component
		Exclude             []ReferableResolvedPath
		ExcludeFilesMatcher []ReferableRegExp
		Integrity           Integrity
	}

	Allow struct {
		DepOnAnyVendor ReferableBool
	}

	Component struct {
		Name           ReferableString
		LocalPathMask  ReferableString
		ResolvedPaths  []ReferableResolvedPath
		AllowedImports []ReferableResolvedPath
		MayDependOn    []ReferableString
		CanUse         []ReferableString
		SpecialFlags   SpecialFlags
	}

	SpecialFlags struct {
		AllowAllProjectDeps ReferableBool
		AllowAllVendorDeps  ReferableBool
	}

	Integrity struct {
		DocumentNotices []Notice
		Suggestions     []Notice
	}

	Notice struct {
		Notice error
		Ref    models.Reference
	}
)
