package speca

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	Spec struct {
		RootDirectory       Referable[string]
		WorkingDirectory    Referable[string]
		ModuleName          Referable[string]
		Allow               Allow
		Components          []Component
		Exclude             []Referable[models.ResolvedPath]
		ExcludeFilesMatcher []Referable[*regexp.Regexp]
		Integrity           Integrity
	}

	Allow struct {
		DepOnAnyVendor Referable[bool]
	}

	Component struct {
		Name                  Referable[string]
		ResolvedPaths         []Referable[models.ResolvedPath]
		AllowedProjectImports []Referable[models.ResolvedPath]
		AllowedVendorGlobs    []Referable[models.Glob]
		MayDependOn           []Referable[string]
		CanUse                []Referable[string]
		SpecialFlags          SpecialFlags
	}

	SpecialFlags struct {
		AllowAllProjectDeps Referable[bool]
		AllowAllVendorDeps  Referable[bool]
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
