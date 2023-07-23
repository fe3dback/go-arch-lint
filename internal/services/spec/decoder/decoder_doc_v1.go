package decoder

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type (
	// ArchV1 initial config schema
	ArchV1 struct {
		FVersion            ref[int]                                    `json:"version"`
		FAllow              ArchV1Allow                                 `json:"allow"`
		FExclude            []ref[string]                               `json:"exclude"`
		FExcludeFilesRegExp []ref[string]                               `json:"excludeFiles"`
		FVendors            map[spec.VendorName]ref[ArchV1Vendor]       `json:"vendors"`
		FCommonVendors      []ref[string]                               `json:"commonVendors"`
		FComponents         map[spec.ComponentName]ref[ArchV1Component] `json:"components"`
		FCommonComponents   []ref[string]                               `json:"commonComponents"`
		FDependencies       map[spec.ComponentName]ref[ArchV1Rule]      `json:"deps"`
	}

	ArchV1Allow struct {
		FDepOnAnyVendor ref[bool] `json:"depOnAnyVendor"`
	}

	ArchV1Vendor struct {
		FImportPath string `json:"in"`
	}

	ArchV1Component struct {
		FLocalPath string `json:"in"`
	}

	ArchV1Rule struct {
		FMayDependOn    []ref[string] `json:"mayDependOn"`
		FCanUse         []ref[string] `json:"canUse"`
		FAnyProjectDeps ref[bool]     `json:"anyProjectDeps"`
		FAnyVendorDeps  ref[bool]     `json:"anyVendorDeps"`
	}
)

func (a *ArchV1) Version() common.Referable[int] {
	return castRef(a.FVersion)
}

func (a *ArchV1) WorkingDirectory() common.Referable[string] {
	return common.NewEmptyReferable("./")
}

func (a *ArchV1) Options() spec.Options {
	return a.FAllow
}

func (a *ArchV1) ExcludedDirectories() []common.Referable[string] {
	return castRefList(a.FExclude)
}

func (a *ArchV1) ExcludedFilesRegExp() []common.Referable[string] {
	return castRefList(a.FExcludeFilesRegExp)
}

func (a *ArchV1) Vendors() spec.Vendors {
	casted := make(spec.Vendors, len(a.FVendors))
	for name, vendor := range a.FVendors {
		casted[name] = common.NewReferable(spec.Vendor(vendor.Value), vendor.Reference)
	}

	return casted
}

func (a *ArchV1) CommonVendors() []common.Referable[string] {
	return castRefList(a.FCommonVendors)
}

func (a *ArchV1) Components() spec.Components {
	casted := make(spec.Components, len(a.FComponents))
	for name, cmp := range a.FComponents {
		casted[name] = common.NewReferable(spec.Component(cmp.Value), cmp.Reference)
	}

	return casted
}

func (a *ArchV1) CommonComponents() []common.Referable[string] {
	return castRefList(a.FCommonComponents)
}

func (a *ArchV1) Dependencies() spec.Dependencies {
	casted := make(spec.Dependencies, len(a.FDependencies))
	for name, dep := range a.FDependencies {
		casted[name] = common.NewReferable(spec.DependencyRule(dep.Value), dep.Reference)
	}

	return casted
}

// --

func (a ArchV1Allow) IsDependOnAnyVendor() common.Referable[bool] {
	return castRef(a.FDepOnAnyVendor)
}

func (a ArchV1Allow) DeepScan() common.Referable[bool] {
	return common.NewEmptyReferable(false)
}

// --

func (a ArchV1Vendor) ImportPaths() []models.Glob {
	return []models.Glob{models.Glob(a.FImportPath)}
}

// --

func (a ArchV1Component) RelativePaths() []models.Glob {
	return []models.Glob{models.Glob(a.FLocalPath)}
}

// --

func (a ArchV1Rule) MayDependOn() []common.Referable[string] {
	return castRefList(a.FMayDependOn)
}

func (a ArchV1Rule) CanUse() []common.Referable[string] {
	return castRefList(a.FCanUse)
}

func (a ArchV1Rule) AnyProjectDeps() common.Referable[bool] {
	return castRef(a.FAnyProjectDeps)
}

func (a ArchV1Rule) AnyVendorDeps() common.Referable[bool] {
	return castRef(a.FAnyVendorDeps)
}

func (a ArchV1Rule) DeepScan() common.Referable[bool] {
	return common.NewEmptyReferable(false)
}
