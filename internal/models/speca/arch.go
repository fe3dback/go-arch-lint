package speca

type (
	Spec struct {
		RootDirectory       ReferableString
		ModuleName          ReferableString
		Allow               Allow
		Components          []Component
		Exclude             []ReferableResolvedPath
		ExcludeFilesMatcher []ReferableRegExp
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
)
