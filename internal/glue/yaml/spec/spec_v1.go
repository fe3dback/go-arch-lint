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
		internalVendors            archV1InternalVendors
		internalComponents         archV1InternalComponents
		internalExclude            archV1InternalExclude
		internalExcludeFilesRegExp archV1InternalExcludeFilesRegExp
		internalCommonVendors      archV1InternalCommonVendors
		internalCommonComponents   archV1InternalCommonComponents
		internalDependencies       archV1InternalDependencies

		V1Version            speca.ReferableInt                     `yaml:"version" json:"version"`
		V1Allow              ArchV1Allow                            `yaml:"allow" json:"allow"`
		V1Vendors            map[arch.VendorName]ArchV1Vendor       `yaml:"vendors" json:"vendors"`
		V1Exclude            []speca.ReferableString                `yaml:"exclude" json:"exclude"`
		V1ExcludeFilesRegExp []speca.ReferableString                `yaml:"excludeFiles" json:"excludeFiles"`
		V1Components         map[arch.ComponentName]ArchV1Component `yaml:"components" json:"components"`
		V1Dependencies       map[arch.ComponentName]ArchV1Rules     `yaml:"deps" json:"deps"`
		V1CommonComponents   []speca.ReferableString                `yaml:"commonComponents" json:"commonComponents"`
		V1CommonVendors      []speca.ReferableString                `yaml:"commonVendors" json:"commonVendors"`
	}

	ArchV1Allow struct {
		reference models.Reference

		V1DepOnAnyVendor speca.ReferableBool `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
	}

	ArchV1Vendor struct {
		reference models.Reference

		V1ImportPath speca.ReferableString `yaml:"in" json:"in"`
	}

	ArchV1Component struct {
		reference models.Reference

		V1LocalPath speca.ReferableString `yaml:"in" json:"in"`
	}

	ArchV1Rules struct {
		reference models.Reference

		V1MayDependOn    []speca.ReferableString `yaml:"mayDependOn" json:"mayDependOn"`
		V1CanUse         []speca.ReferableString `yaml:"canUse" json:"canUse"`
		V1AnyProjectDeps speca.ReferableBool     `yaml:"anyProjectDeps" json:"anyProjectDeps"`
		V1AnyVendorDeps  speca.ReferableBool     `yaml:"anyVendorDeps" json:"anyVendorDeps"`
	}

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
	return doc.V1Version
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
	// Version
	doc.V1Version = speca.NewReferableInt(
		doc.V1Version.Value(),
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
			item.Value(),
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
			item.Value(),
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
			item.Value(),
			resolver.Resolve(fmt.Sprintf("$.commonComponents[%d]", ind)),
		)
	}
	doc.internalCommonComponents = archV1InternalCommonComponents{
		reference: resolver.Resolve("$.commonComponents"),
		data:      commonComponents,
	}

	// CommonVendors // todo
	commonVendors := make([]speca.ReferableString, len(doc.V1CommonVendors))
	for ind, item := range doc.V1CommonVendors {
		commonVendors[ind] = speca.NewReferableString(
			item.Value(),
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
	return opt.V1DepOnAnyVendor
}

func (opt ArchV1Allow) applyReferences(resolver YamlSourceCodeReferenceResolver) ArchV1Allow {
	opt.reference = resolver.Resolve("$.allow")

	opt.V1DepOnAnyVendor = speca.NewReferableBool(
		opt.V1DepOnAnyVendor.Value(),
		resolver.Resolve("$.allow.depOnAnyVendor"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV1Vendor) Reference() models.Reference {
	return v.reference
}

func (v ArchV1Vendor) ImportPath() speca.ReferableString {
	return v.V1ImportPath
}

func (v ArchV1Vendor) applyReferences(name arch.VendorName, resolver YamlSourceCodeReferenceResolver) ArchV1Vendor {
	v.reference = resolver.Resolve(fmt.Sprintf("$.vendors.%s", name))
	v.V1ImportPath = speca.NewReferableString(
		v.V1ImportPath.Value(),
		resolver.Resolve(fmt.Sprintf("$.vendors.%s.in", name)),
	)

	return v
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (c ArchV1Component) Reference() models.Reference {
	return c.reference
}

func (c ArchV1Component) LocalPath() speca.ReferableString {
	return c.V1LocalPath
}

func (c ArchV1Component) applyReferences(name arch.ComponentName, resolver YamlSourceCodeReferenceResolver) ArchV1Component {
	c.reference = resolver.Resolve(fmt.Sprintf("$.components.%s", name))
	c.V1LocalPath = speca.NewReferableString(
		c.V1LocalPath.Value(),
		resolver.Resolve(fmt.Sprintf("$.components.%s.in", name)),
	)

	return c
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (rule ArchV1Rules) Reference() models.Reference {
	return rule.reference
}

func (rule ArchV1Rules) MayDependOn() []speca.ReferableString {
	return rule.V1MayDependOn
}

func (rule ArchV1Rules) CanUse() []speca.ReferableString {
	return rule.V1CanUse
}

func (rule ArchV1Rules) AnyProjectDeps() speca.ReferableBool {
	return rule.V1AnyProjectDeps
}

func (rule ArchV1Rules) AnyVendorDeps() speca.ReferableBool {
	return rule.V1AnyVendorDeps
}

func (rule ArchV1Rules) applyReferences(name arch.ComponentName, resolver YamlSourceCodeReferenceResolver) ArchV1Rules {
	rule.reference = resolver.Resolve(fmt.Sprintf("$.deps.%s", name))

	// --
	rule.V1AnyVendorDeps = speca.NewReferableBool(
		rule.V1AnyVendorDeps.Value(),
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyVendorDeps", name)),
	)

	// --
	rule.V1AnyProjectDeps = speca.NewReferableBool(
		rule.V1AnyProjectDeps.Value(),
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyProjectDeps", name)),
	)

	// --
	canUse := make([]speca.ReferableString, len(rule.V1CanUse))
	for ind, item := range rule.V1CanUse {
		canUse[ind] = speca.NewReferableString(
			item.Value(),
			resolver.Resolve(fmt.Sprintf("$.deps.%s.canUse[%d]", name, ind)),
		)
	}
	rule.V1CanUse = canUse

	// --
	mayDependOn := make([]speca.ReferableString, len(rule.V1MayDependOn))
	for ind, item := range rule.V1MayDependOn {
		mayDependOn[ind] = speca.NewReferableString(
			item.Value(),
			resolver.Resolve(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, ind)),
		)
	}
	rule.V1MayDependOn = mayDependOn

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
