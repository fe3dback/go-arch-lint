package arch

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type (
	Spec struct {
		RootDirectory       common.Referable[string]
		WorkingDirectory    common.Referable[string]
		ModuleName          common.Referable[string]
		Allow               Allow
		Components          []Component
		Exclude             []common.Referable[models.ResolvedPath]
		ExcludeFilesMatcher []common.Referable[*regexp.Regexp]
		Integrity           Integrity
	}

	Allow struct {
		DepOnAnyVendor common.Referable[bool]
		DeepScan       common.Referable[bool]
	}

	Component struct {
		Name                  common.Referable[string]
		DeepScan              common.Referable[bool]
		ResolvedPaths         []common.Referable[models.ResolvedPath]
		AllowedProjectImports []common.Referable[models.ResolvedPath]
		AllowedVendorGlobs    []common.Referable[models.Glob]
		MayDependOn           []common.Referable[string]
		CanUse                []common.Referable[string]
		SpecialFlags          SpecialFlags
	}

	SpecialFlags struct {
		AllowAllProjectDeps common.Referable[bool]
		AllowAllVendorDeps  common.Referable[bool]
	}

	Integrity struct {
		DocumentNotices []Notice
		Suggestions     []Notice
	}

	Notice struct {
		Notice error
		Ref    common.Reference
	}
)
