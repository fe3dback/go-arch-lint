package decoder

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type (
	// ArchV2 changes since ArchV1:
	// - added global "workdir" option
	// - vendor/component "in" now accept two types: (string, []string), instead of (string)
	ArchV2 struct {
		FVersion            ref[int]                                    `json:"version"`
		FWorkDir            ref[string]                                 `json:"workdir"`
		FAllow              ArchV2Allow                                 `json:"allow"`
		FExclude            []ref[string]                               `json:"exclude"`
		FExcludeFilesRegExp []ref[string]                               `json:"excludeFiles"`
		FVendors            map[spec.VendorName]ref[ArchV2Vendor]       `json:"vendors"`
		FCommonVendors      []ref[string]                               `json:"commonVendors"`
		FComponents         map[spec.ComponentName]ref[ArchV2Component] `json:"components"`
		FCommonComponents   []ref[string]                               `json:"commonComponents"`
		FDependencies       map[spec.ComponentName]ref[ArchV2Rule]      `json:"deps"`
	}

	ArchV2Allow struct {
		FDepOnAnyVendor ref[bool] `json:"depOnAnyVendor"`
	}

	ArchV2Vendor struct {
		FImportPaths stringList `json:"in"`
	}

	ArchV2Component struct {
		FLocalPaths stringList `json:"in"`
	}

	ArchV2Rule struct {
		FMayDependOn    []ref[string] `json:"mayDependOn"`
		FCanUse         []ref[string] `json:"canUse"`
		FAnyProjectDeps ref[bool]     `json:"anyProjectDeps"`
		FAnyVendorDeps  ref[bool]     `json:"anyVendorDeps"`
	}
)

func (a *ArchV2) postSetup() {}

func (a *ArchV2) Version() common.Referable[int] {
	return castRef(a.FVersion)
}

func (a *ArchV2) WorkingDirectory() common.Referable[string] {
	// fallback from version 1
	actualWorkDirectory := "./"

	if a.FWorkDir.ref.Value != "" {
		actualWorkDirectory = a.FWorkDir.ref.Value
	}

	return common.NewReferable(actualWorkDirectory, a.FWorkDir.ref.Reference)
}

func (a *ArchV2) Options() spec.Options {
	return a.FAllow
}

func (a *ArchV2) ExcludedDirectories() []common.Referable[string] {
	return castRefList(a.FExclude)
}

func (a *ArchV2) ExcludedFilesRegExp() []common.Referable[string] {
	return castRefList(a.FExcludeFilesRegExp)
}

func (a *ArchV2) Vendors() spec.Vendors {
	casted := make(spec.Vendors, len(a.FVendors))
	for name, vendor := range a.FVendors {
		casted[name] = common.NewReferable(spec.Vendor(vendor.ref.Value), vendor.ref.Reference)
	}

	return casted
}

func (a *ArchV2) CommonVendors() []common.Referable[string] {
	return castRefList(a.FCommonVendors)
}

func (a *ArchV2) Components() spec.Components {
	casted := make(spec.Components, len(a.FComponents))
	for name, cmp := range a.FComponents {
		casted[name] = common.NewReferable(spec.Component(cmp.ref.Value), cmp.ref.Reference)
	}

	return casted
}

func (a *ArchV2) CommonComponents() []common.Referable[string] {
	return castRefList(a.FCommonComponents)
}

func (a *ArchV2) Dependencies() spec.Dependencies {
	casted := make(spec.Dependencies, len(a.FDependencies))
	for name, dep := range a.FDependencies {
		casted[name] = common.NewReferable(spec.DependencyRule(dep.ref.Value), dep.ref.Reference)
	}

	return casted
}

// --

func (a ArchV2Allow) IsDependOnAnyVendor() common.Referable[bool] {
	return castRef(a.FDepOnAnyVendor)
}

func (a ArchV2Allow) DeepScan() common.Referable[bool] {
	return common.NewEmptyReferable(false)
}

// --

func (a ArchV2Vendor) ImportPaths() []models.Glob {
	casted := make([]models.Glob, 0, len(a.FImportPaths))

	for _, path := range a.FImportPaths {
		casted = append(casted, models.Glob(path))
	}

	return casted
}

// --

func (a ArchV2Component) RelativePaths() []models.Glob {
	casted := make([]models.Glob, 0, len(a.FLocalPaths))

	for _, path := range a.FLocalPaths {
		casted = append(casted, models.Glob(path))
	}

	return casted
}

// --

func (a ArchV2Rule) MayDependOn() []common.Referable[string] {
	return castRefList(a.FMayDependOn)
}

func (a ArchV2Rule) CanUse() []common.Referable[string] {
	return castRefList(a.FCanUse)
}

func (a ArchV2Rule) AnyProjectDeps() common.Referable[bool] {
	return castRef(a.FAnyProjectDeps)
}

func (a ArchV2Rule) AnyVendorDeps() common.Referable[bool] {
	return castRef(a.FAnyVendorDeps)
}

func (a ArchV2Rule) DeepScan() common.Referable[bool] {
	return common.NewEmptyReferable(false)
}
