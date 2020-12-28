package spec

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ArchV1Document struct {
		reference                  models.Reference
		internalVersion            speca.ReferableInt
		internalVendors            archV1InternalVendors
		internalExclude            archV1InternalExclude
		internalExcludeFilesRegExp archV1InternalExcludeFilesRegExp
		internalComponents         archV1InternalComponents
		internalDependencies       archV1InternalDependencies
		internalCommonComponents   archV1InternalCommonComponents
		internalCommonVendors      archV1InternalCommonVendors

		V1Version            int                                    `yaml:"version" json:"version"`
		V1Allow              ArchV1Allow                            `yaml:"allow" json:"allow"`
		V1Vendors            map[arch.VendorName]ArchV1Vendor       `yaml:"vendors" json:"vendors"`
		V1Exclude            []string                               `yaml:"exclude" json:"exclude"`
		V1ExcludeFilesRegExp []string                               `yaml:"excludeFiles" json:"excludeFiles"`
		V1Components         map[arch.ComponentName]ArchV1Component `yaml:"components" json:"components"`
		V1Dependencies       map[arch.ComponentName]ArchV1Rules     `yaml:"deps" json:"deps"`
		V1CommonComponents   []string                               `yaml:"commonComponents" json:"commonComponents"`
		V1CommonVendors      []string                               `yaml:"commonVendors" json:"commonVendors"`
	}

	ArchV1Allow struct {
		reference              models.Reference
		internalDepOnAnyVendor speca.ReferableBool

		V1DepOnAnyVendor bool `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
	}

	ArchV1Vendor struct {
		reference          models.Reference
		internalImportPath speca.ReferableString

		V1ImportPath string `yaml:"in" json:"in"`
	}

	ArchV1Component struct {
		reference         models.Reference
		internalLocalPath speca.ReferableString

		V1LocalPath string `yaml:"in" json:"in"`
	}

	ArchV1Rules struct {
		reference              models.Reference
		internalMayDependOn    []speca.ReferableString
		internalCanUse         []speca.ReferableString
		internalAnyProjectDeps speca.ReferableBool
		internalAnyVendorDeps  speca.ReferableBool

		V1MayDependOn    []string `yaml:"mayDependOn" json:"mayDependOn"`
		V1CanUse         []string `yaml:"canUse" json:"canUse"`
		V1AnyProjectDeps bool     `yaml:"anyProjectDeps" json:"anyProjectDeps"`
		V1AnyVendorDeps  bool     `yaml:"anyVendorDeps" json:"anyVendorDeps"`
	}
)

type (
	archV1InternalVendors struct {
		reference models.Reference
		data      map[arch.VendorName]ArchV1Vendor
	}

	archV1InternalComponents struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV1Component
	}

	archV1InternalExclude struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalExcludeFilesRegExp struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalCommonVendors struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalCommonComponents struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalDependencies struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV1Rules
	}
)

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (doc ArchV1Document) Reference() models.Reference {
	return doc.reference
}

func (doc ArchV1Document) Version() speca.ReferableInt {
	return doc.internalVersion
}

func (doc ArchV1Document) WorkingDirectory() speca.ReferableString {
	return speca.NewReferableString(
		"./",
		speca.NewEmptyReference(),
	)
}

func (doc ArchV1Document) Options() arch.Options {
	return doc.V1Allow
}

func (doc ArchV1Document) ExcludedDirectories() arch.ExcludedDirectories {
	return doc.internalExclude
}

func (doc ArchV1Document) ExcludedFilesRegExp() arch.ExcludedFilesRegExp {
	return doc.internalExcludeFilesRegExp
}

func (doc ArchV1Document) Vendors() arch.Vendors {
	return doc.internalVendors
}

func (doc ArchV1Document) Components() arch.Components {
	return doc.internalComponents
}

func (doc ArchV1Document) CommonComponents() arch.CommonComponents {
	return doc.internalCommonComponents
}

func (doc ArchV1Document) CommonVendors() arch.CommonVendors {
	return doc.internalCommonVendors
}

func (doc ArchV1Document) Dependencies() arch.Dependencies {
	return doc.internalDependencies
}

func (doc ArchV1Document) applyReferences(resolver YamlSourceCodeReferenceResolver) ArchV1Document {
	doc.reference = resolver.Resolve("$.version")

	// Version
	doc.internalVersion = speca.NewReferableInt(
		doc.V1Version,
		resolver.Resolve("$.version"),
	)

	// Allow
	doc.V1Allow = doc.V1Allow.applyReferences(resolver)

	// Vendors
	vendors := make(map[string]ArchV1Vendor)
	for name, vendor := range doc.V1Vendors {
		vendors[name] = vendor.applyReferences(name, resolver)
	}
	doc.internalVendors = archV1InternalVendors{
		reference: resolver.Resolve("$.vendors"),
		data:      vendors,
	}

	// Exclude
	excludedDirectories := make([]speca.ReferableString, len(doc.V1Exclude))
	for ind, item := range doc.V1Exclude {
		excludedDirectories[ind] = speca.NewReferableString(
			item,
			resolver.Resolve(fmt.Sprintf("$.exclude[%d]", ind)),
		)
	}

	doc.internalExclude = archV1InternalExclude{
		reference: resolver.Resolve("$.exclude"),
		data:      excludedDirectories,
	}

	// ExcludeFilesRegExp
	excludedFiles := make([]speca.ReferableString, len(doc.V1ExcludeFilesRegExp))
	for ind, item := range doc.V1ExcludeFilesRegExp {
		excludedFiles[ind] = speca.NewReferableString(
			item,
			resolver.Resolve(fmt.Sprintf("$.excludeFiles[%d]", ind)),
		)
	}

	doc.internalExcludeFilesRegExp = archV1InternalExcludeFilesRegExp{
		reference: resolver.Resolve("$.excludeFiles"),
		data:      excludedFiles,
	}

	// Components
	components := make(map[string]ArchV1Component)
	for name, component := range doc.V1Components {
		components[name] = component.applyReferences(name, resolver)
	}
	doc.internalComponents = archV1InternalComponents{
		reference: resolver.Resolve("$.components"),
		data:      components,
	}

	// Dependencies
	dependencies := make(map[string]ArchV1Rules)
	for name, rules := range doc.V1Dependencies {
		dependencies[name] = rules.applyReferences(name, resolver)
	}
	doc.internalDependencies = archV1InternalDependencies{
		reference: resolver.Resolve("$.deps"),
		data:      dependencies,
	}

	// CommonComponents
	commonComponents := make([]speca.ReferableString, len(doc.V1CommonComponents))
	for ind, item := range doc.V1CommonComponents {
		commonComponents[ind] = speca.NewReferableString(
			item,
			resolver.Resolve(fmt.Sprintf("$.commonComponents[%d]", ind)),
		)
	}
	doc.internalCommonComponents = archV1InternalCommonComponents{
		reference: resolver.Resolve("$.commonComponents"),
		data:      commonComponents,
	}

	// CommonVendors
	commonVendors := make([]speca.ReferableString, len(doc.V1CommonVendors))
	for ind, item := range doc.V1CommonVendors {
		commonVendors[ind] = speca.NewReferableString(
			item,
			resolver.Resolve(fmt.Sprintf("$.commonVendors[%d]", ind)),
		)
	}
	doc.internalCommonVendors = archV1InternalCommonVendors{
		reference: resolver.Resolve("$.commonVendors"),
		data:      commonVendors,
	}

	return doc
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (opt ArchV1Allow) Reference() models.Reference {
	return opt.reference
}

func (opt ArchV1Allow) IsDependOnAnyVendor() speca.ReferableBool {
	return opt.internalDepOnAnyVendor
}

func (opt ArchV1Allow) applyReferences(resolver YamlSourceCodeReferenceResolver) ArchV1Allow {
	opt.reference = resolver.Resolve("$.allow")

	opt.internalDepOnAnyVendor = speca.NewReferableBool(
		opt.V1DepOnAnyVendor,
		resolver.Resolve("$.allow.depOnAnyVendor"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV1Vendor) Reference() models.Reference {
	return v.reference
}

func (v ArchV1Vendor) ImportPath() speca.ReferableString {
	return v.internalImportPath
}

func (v ArchV1Vendor) applyReferences(name arch.VendorName, resolver YamlSourceCodeReferenceResolver) ArchV1Vendor {
	v.reference = resolver.Resolve(fmt.Sprintf("$.vendors.%s", name))
	v.internalImportPath = speca.NewReferableString(
		v.V1ImportPath,
		resolver.Resolve(fmt.Sprintf("$.vendors.%s.in", name)),
	)

	return v
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (c ArchV1Component) Reference() models.Reference {
	return c.reference
}

func (c ArchV1Component) LocalPath() speca.ReferableString {
	return c.internalLocalPath
}

func (c ArchV1Component) applyReferences(name arch.ComponentName, resolver YamlSourceCodeReferenceResolver) ArchV1Component {
	c.reference = resolver.Resolve(fmt.Sprintf("$.components.%s", name))
	c.internalLocalPath = speca.NewReferableString(
		c.V1LocalPath,
		resolver.Resolve(fmt.Sprintf("$.components.%s.in", name)),
	)

	return c
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (rule ArchV1Rules) Reference() models.Reference {
	return rule.reference
}

func (rule ArchV1Rules) MayDependOn() []speca.ReferableString {
	return rule.internalMayDependOn
}

func (rule ArchV1Rules) CanUse() []speca.ReferableString {
	return rule.internalCanUse
}

func (rule ArchV1Rules) AnyProjectDeps() speca.ReferableBool {
	return rule.internalAnyProjectDeps
}

func (rule ArchV1Rules) AnyVendorDeps() speca.ReferableBool {
	return rule.internalAnyVendorDeps
}

func (rule ArchV1Rules) applyReferences(name arch.ComponentName, resolver YamlSourceCodeReferenceResolver) ArchV1Rules {
	rule.reference = resolver.Resolve(fmt.Sprintf("$.deps.%s", name))

	// --
	rule.internalAnyVendorDeps = speca.NewReferableBool(
		rule.V1AnyVendorDeps,
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyVendorDeps", name)),
	)

	// --
	rule.internalAnyProjectDeps = speca.NewReferableBool(
		rule.V1AnyProjectDeps,
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyProjectDeps", name)),
	)

	// --
	canUse := make([]speca.ReferableString, len(rule.V1CanUse))
	for ind, item := range rule.V1CanUse {
		canUse[ind] = speca.NewReferableString(
			item,
			resolver.Resolve(fmt.Sprintf("$.deps.%s.canUse[%d]", name, ind)),
		)
	}
	rule.internalCanUse = canUse

	// --
	mayDependOn := make([]speca.ReferableString, len(rule.V1MayDependOn))
	for ind, item := range rule.V1MayDependOn {
		mayDependOn[ind] = speca.NewReferableString(
			item,
			resolver.Resolve(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, ind)),
		)
	}
	rule.internalMayDependOn = mayDependOn

	// --
	return rule
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (a archV1InternalDependencies) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalDependencies) Map() map[arch.ComponentName]arch.DependencyRule {
	res := make(map[arch.ComponentName]arch.DependencyRule)
	for name, rules := range a.data {
		res[name] = rules
	}
	return res
}

func (a archV1InternalCommonComponents) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalCommonComponents) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalCommonVendors) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalCommonVendors) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalExcludeFilesRegExp) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalExcludeFilesRegExp) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalExclude) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalExclude) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalComponents) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalComponents) Map() map[arch.ComponentName]arch.Component {
	res := make(map[arch.ComponentName]arch.Component)
	for name, component := range a.data {
		res[name] = component
	}
	return res
}

func (a archV1InternalVendors) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalVendors) Map() map[arch.VendorName]arch.Vendor {
	res := make(map[arch.VendorName]arch.Vendor)
	for name, vendor := range a.data {
		res[name] = vendor
	}
	return res
}
