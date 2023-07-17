package decoder

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type (
	ArchV3 struct {
		FVersion            ref[int]                                    `json:"version"`
		FWorkDir            ref[string]                                 `json:"workdir"`
		FAllow              ArchV3Allow                                 `json:"allow"`
		FExclude            []ref[string]                               `json:"exclude"`
		FExcludeFilesRegExp []ref[string]                               `json:"excludeFiles"`
		FVendors            map[spec.VendorName]ref[ArchV3Vendor]       `json:"vendors"`
		FCommonVendors      []ref[string]                               `json:"commonVendors"`
		FComponents         map[spec.ComponentName]ref[ArchV3Component] `json:"components"`
		FCommonComponents   []ref[string]                               `json:"commonComponents"`
		FDependencies       map[spec.ComponentName]ref[ArchV3Rule]      `json:"deps"`
	}

	ArchV3Allow struct {
		FDepOnAnyVendor ref[bool]  `json:"depOnAnyVendor"`
		FDeepScan       ref[*bool] `json:"deepScan"`
	}

	ArchV3Vendor struct {
		FImportPaths stringList `json:"in"`
	}

	ArchV3Component struct {
		FLocalPaths stringList `json:"in"`
	}

	ArchV3Rule struct {
		FMayDependOn    []ref[string] `json:"mayDependOn"`
		FCanUse         []ref[string] `json:"canUse"`
		FAnyProjectDeps ref[bool]     `json:"anyProjectDeps"`
		FAnyVendorDeps  ref[bool]     `json:"anyVendorDeps"`
		FDeepScan       ref[*bool]    `json:"deepScan"`
	}
)

func (a *ArchV3) Version() common.Referable[int] {
	return castRef(a.FVersion)
}

func (a *ArchV3) WorkingDirectory() common.Referable[string] {
	// fallback from version 1
	actualWorkDirectory := "./"

	if a.FWorkDir.Value != "" {
		actualWorkDirectory = a.FWorkDir.Value
	}

	return common.NewReferable(actualWorkDirectory, a.FWorkDir.Reference)
}

func (a *ArchV3) Options() spec.Options {
	return a.FAllow
}

func (a *ArchV3) ExcludedDirectories() []common.Referable[string] {
	return castRefList(a.FExclude)
}

func (a *ArchV3) ExcludedFilesRegExp() []common.Referable[string] {
	return castRefList(a.FExcludeFilesRegExp)
}

func (a *ArchV3) Vendors() spec.Vendors {
	casted := make(spec.Vendors, len(a.FVendors))
	for name, vendor := range a.FVendors {
		casted[name] = common.NewReferable(spec.Vendor(vendor.Value), vendor.Reference)
	}

	return casted
}

func (a *ArchV3) CommonVendors() []common.Referable[string] {
	return castRefList(a.FCommonVendors)
}

func (a *ArchV3) Components() spec.Components {
	casted := make(spec.Components, len(a.FComponents))
	for name, cmp := range a.FComponents {
		casted[name] = common.NewReferable(spec.Component(cmp.Value), cmp.Reference)
	}

	return casted
}

func (a *ArchV3) CommonComponents() []common.Referable[string] {
	return castRefList(a.FCommonComponents)
}

func (a *ArchV3) Dependencies() spec.Dependencies {
	casted := make(spec.Dependencies, len(a.FDependencies))
	for name, dep := range a.FDependencies {
		casted[name] = common.NewReferable(spec.DependencyRule(dep.Value), dep.Reference)
	}

	return casted
}

// --

func (a ArchV3Allow) IsDependOnAnyVendor() common.Referable[bool] {
	return castRef(a.FDepOnAnyVendor)
}

func (a ArchV3Allow) DeepScan() common.Referable[bool] {
	deepScan := false
	if a.FDeepScan.Value == nil {
		// be default it`s on from V3+
		deepScan = true
	} else {
		deepScan = *a.FDeepScan.Value
	}

	return common.NewReferable(deepScan, a.FDeepScan.Reference)
}

// --

func (a ArchV3Vendor) ImportPaths() []models.Glob {
	casted := make([]models.Glob, 0, len(a.FImportPaths))

	for _, path := range a.FImportPaths {
		casted = append(casted, models.Glob(path))
	}

	return casted
}

// --

func (a ArchV3Component) RelativePaths() []models.Glob {
	casted := make([]models.Glob, 0, len(a.FLocalPaths))

	for _, path := range a.FLocalPaths {
		casted = append(casted, models.Glob(path))
	}

	return casted
}

// --

func (a ArchV3Rule) MayDependOn() []common.Referable[string] {
	return castRefList(a.FMayDependOn)
}

func (a ArchV3Rule) CanUse() []common.Referable[string] {
	return castRefList(a.FCanUse)
}

func (a ArchV3Rule) AnyProjectDeps() common.Referable[bool] {
	return castRef(a.FAnyProjectDeps)
}

func (a ArchV3Rule) AnyVendorDeps() common.Referable[bool] {
	return castRef(a.FAnyVendorDeps)
}

func (a ArchV3Rule) DeepScan() common.Referable[bool] {
	deepScan := false
	if a.FDeepScan.Value == nil {
		// be default it`s on from V3+
		deepScan = true
	} else {
		deepScan = *a.FDeepScan.Value
	}

	return common.NewReferable(deepScan, a.FDeepScan.Reference)
}
