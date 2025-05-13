package decoder

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type (
	// ArchV3 changes since ArchV2:
	// - added deepScan option in allow and deps rules
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
		FDepOnAnyVendor           ref[bool] `json:"depOnAnyVendor"`
		FDeepScan                 ref[bool] `json:"deepScan"`
		FIgnoreNotFoundComponents ref[bool] `json:"ignoreNotFoundComponents"`
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
		FDeepScan       ref[bool]     `json:"deepScan"`
	}
)

func (a *ArchV3) postSetup() {
	// deep scan nesting (global settings -> local settings)
	for depName := range a.FDependencies {
		localDeepScan := a.FDependencies[depName].ref.Value.FDeepScan

		if !localDeepScan.defined {
			dep := a.FDependencies[depName]
			dep.ref.Value.FDeepScan = ref[bool]{
				defined: true,
				ref:     a.FAllow.DeepScan(),
			}

			a.FDependencies[depName] = dep
		}
	}
}

func (a *ArchV3) Version() common.Referable[int] {
	return castRef(a.FVersion)
}

func (a *ArchV3) WorkingDirectory() common.Referable[string] {
	// fallback from version 1
	actualWorkDirectory := "./"

	if a.FWorkDir.ref.Value != "" {
		actualWorkDirectory = a.FWorkDir.ref.Value
	}

	return common.NewReferable(actualWorkDirectory, a.FWorkDir.ref.Reference)
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
		casted[name] = common.NewReferable(spec.Vendor(vendor.ref.Value), vendor.ref.Reference)
	}

	return casted
}

func (a *ArchV3) CommonVendors() []common.Referable[string] {
	return castRefList(a.FCommonVendors)
}

func (a *ArchV3) Components() spec.Components {
	casted := make(spec.Components, len(a.FComponents))
	for name, cmp := range a.FComponents {
		casted[name] = common.NewReferable(spec.Component(cmp.ref.Value), cmp.ref.Reference)
	}

	return casted
}

func (a *ArchV3) CommonComponents() []common.Referable[string] {
	return castRefList(a.FCommonComponents)
}

func (a *ArchV3) Dependencies() spec.Dependencies {
	casted := make(spec.Dependencies, len(a.FDependencies))
	for name, dep := range a.FDependencies {
		casted[name] = common.NewReferable(spec.DependencyRule(dep.ref.Value), dep.ref.Reference)
	}

	return casted
}

// --

func (a ArchV3Allow) IsDependOnAnyVendor() common.Referable[bool] {
	return castRef(a.FDepOnAnyVendor)
}

func (a ArchV3Allow) DeepScan() common.Referable[bool] {
	if a.FDeepScan.defined {
		return a.FDeepScan.ref
	}

	// be default it`s on from V3+
	return common.NewEmptyReferable(true)
}

func (a ArchV3Allow) IgnoreNotFoundComponents() common.Referable[bool] {
	if a.FIgnoreNotFoundComponents.defined {
		return a.FIgnoreNotFoundComponents.ref
	}

	// disabled by default
	return common.NewEmptyReferable(false)
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
	return a.FDeepScan.ref
}
