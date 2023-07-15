package spec

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ArchV1Document struct {
		filePath speca.Referable[string]

		reference                  models.Reference
		internalVersion            speca.Referable[int]
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
		internalDepOnAnyVendor speca.Referable[bool]

		V1DepOnAnyVendor bool `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
	}

	ArchV1Vendor struct {
		reference          models.Reference
		internalImportPath speca.Referable[models.Glob]

		V1ImportPath string `yaml:"in" json:"in"`
	}

	ArchV1Component struct {
		reference         models.Reference
		internalLocalPath speca.Referable[models.Glob]

		V1LocalPath string `yaml:"in" json:"in"`
	}

	ArchV1Rules struct {
		reference              models.Reference
		internalMayDependOn    []speca.Referable[string]
		internalCanUse         []speca.Referable[string]
		internalAnyProjectDeps speca.Referable[bool]
		internalAnyVendorDeps  speca.Referable[bool]

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
		data      []speca.Referable[string]
	}

	archV1InternalExcludeFilesRegExp struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV1InternalCommonVendors struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV1InternalCommonComponents struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV1InternalDependencies struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV1Rules
	}
)

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (doc ArchV1Document) FilePath() speca.Referable[string] {
	return doc.filePath
}

func (doc ArchV1Document) Reference() models.Reference {
	return doc.reference
}

func (doc ArchV1Document) Version() speca.Referable[int] {
	return doc.internalVersion
}

func (doc ArchV1Document) WorkingDirectory() speca.Referable[string] {
	return speca.NewReferable(
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

func (doc ArchV1Document) applyReferences(resolve yamlDocumentPathResolver) ArchV1Document {
	doc.reference = resolve("$.version")

	// Version
	doc.internalVersion = speca.NewReferable(
		doc.V1Version,
		resolve("$.version"),
	)

	// Allow
	doc.V1Allow = doc.V1Allow.applyReferences(resolve)

	// Vendors
	vendors := make(map[string]ArchV1Vendor)
	for name, vendor := range doc.V1Vendors {
		vendors[name] = vendor.applyReferences(name, resolve)
	}
	doc.internalVendors = archV1InternalVendors{
		reference: resolve("$.vendors"),
		data:      vendors,
	}

	// Exclude
	excludedDirectories := make([]speca.Referable[string], len(doc.V1Exclude))
	for ind, item := range doc.V1Exclude {
		excludedDirectories[ind] = speca.NewReferable(
			item,
			resolve(fmt.Sprintf("$.exclude[%d]", ind)),
		)
	}

	doc.internalExclude = archV1InternalExclude{
		reference: resolve("$.exclude"),
		data:      excludedDirectories,
	}

	// ExcludeFilesRegExp
	excludedFiles := make([]speca.Referable[string], len(doc.V1ExcludeFilesRegExp))
	for ind, item := range doc.V1ExcludeFilesRegExp {
		excludedFiles[ind] = speca.NewReferable(
			item,
			resolve(fmt.Sprintf("$.excludeFiles[%d]", ind)),
		)
	}

	doc.internalExcludeFilesRegExp = archV1InternalExcludeFilesRegExp{
		reference: resolve("$.excludeFiles"),
		data:      excludedFiles,
	}

	// Components
	components := make(map[string]ArchV1Component)
	for name, component := range doc.V1Components {
		components[name] = component.applyReferences(name, resolve)
	}
	doc.internalComponents = archV1InternalComponents{
		reference: resolve("$.components"),
		data:      components,
	}

	// Dependencies
	dependencies := make(map[string]ArchV1Rules)
	for name, rules := range doc.V1Dependencies {
		dependencies[name] = rules.applyReferences(name, resolve)
	}
	doc.internalDependencies = archV1InternalDependencies{
		reference: resolve("$.deps"),
		data:      dependencies,
	}

	// CommonComponents
	commonComponents := make([]speca.Referable[string], len(doc.V1CommonComponents))
	for ind, item := range doc.V1CommonComponents {
		commonComponents[ind] = speca.NewReferable(
			item,
			resolve(fmt.Sprintf("$.commonComponents[%d]", ind)),
		)
	}
	doc.internalCommonComponents = archV1InternalCommonComponents{
		reference: resolve("$.commonComponents"),
		data:      commonComponents,
	}

	// CommonVendors
	commonVendors := make([]speca.Referable[string], len(doc.V1CommonVendors))
	for ind, item := range doc.V1CommonVendors {
		commonVendors[ind] = speca.NewReferable(
			item,
			resolve(fmt.Sprintf("$.commonVendors[%d]", ind)),
		)
	}
	doc.internalCommonVendors = archV1InternalCommonVendors{
		reference: resolve("$.commonVendors"),
		data:      commonVendors,
	}

	return doc
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (opt ArchV1Allow) Reference() models.Reference {
	return opt.reference
}

func (opt ArchV1Allow) IsDependOnAnyVendor() speca.Referable[bool] {
	return opt.internalDepOnAnyVendor
}

func (opt ArchV1Allow) DeepScan() speca.Referable[bool] {
	return speca.NewEmptyReferable(false)
}

func (opt ArchV1Allow) applyReferences(resolver yamlDocumentPathResolver) ArchV1Allow {
	opt.reference = resolver("$.allow")

	opt.internalDepOnAnyVendor = speca.NewReferable(
		opt.V1DepOnAnyVendor,
		resolver("$.allow.depOnAnyVendor"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV1Vendor) Reference() models.Reference {
	return v.reference
}

func (v ArchV1Vendor) ImportPaths() []speca.Referable[models.Glob] {
	return []speca.Referable[models.Glob]{
		v.internalImportPath,
	}
}

func (v ArchV1Vendor) applyReferences(name arch.VendorName, resolve yamlDocumentPathResolver) ArchV1Vendor {
	v.reference = resolve(fmt.Sprintf("$.vendors.%s", name))
	v.internalImportPath = speca.NewReferable(
		models.Glob(v.V1ImportPath),
		resolve(fmt.Sprintf("$.vendors.%s.in", name)),
	)

	return v
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (c ArchV1Component) Reference() models.Reference {
	return c.reference
}

func (c ArchV1Component) RelativePaths() []speca.Referable[models.Glob] {
	return []speca.Referable[models.Glob]{
		c.internalLocalPath,
	}
}

func (c ArchV1Component) applyReferences(name arch.ComponentName, resolve yamlDocumentPathResolver) ArchV1Component {
	c.reference = resolve(fmt.Sprintf("$.components.%s", name))
	c.internalLocalPath = speca.NewReferable(
		models.Glob(c.V1LocalPath),
		resolve(fmt.Sprintf("$.components.%s.in", name)),
	)

	return c
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (rule ArchV1Rules) Reference() models.Reference {
	return rule.reference
}

func (rule ArchV1Rules) MayDependOn() []speca.Referable[string] {
	return rule.internalMayDependOn
}

func (rule ArchV1Rules) CanUse() []speca.Referable[string] {
	return rule.internalCanUse
}

func (rule ArchV1Rules) AnyProjectDeps() speca.Referable[bool] {
	return rule.internalAnyProjectDeps
}

func (rule ArchV1Rules) AnyVendorDeps() speca.Referable[bool] {
	return rule.internalAnyVendorDeps
}

func (rule ArchV1Rules) DeepScan() speca.Referable[bool] {
	return speca.NewEmptyReferable(false)
}

func (rule ArchV1Rules) applyReferences(name arch.ComponentName, resolve yamlDocumentPathResolver) ArchV1Rules {
	rule.reference = resolve(fmt.Sprintf("$.deps.%s", name))

	// --
	rule.internalAnyVendorDeps = speca.NewReferable(
		rule.V1AnyVendorDeps,
		resolve(fmt.Sprintf("$.deps.%s.anyVendorDeps", name)),
	)

	// --
	rule.internalAnyProjectDeps = speca.NewReferable(
		rule.V1AnyProjectDeps,
		resolve(fmt.Sprintf("$.deps.%s.anyProjectDeps", name)),
	)

	// --
	canUse := make([]speca.Referable[string], len(rule.V1CanUse))
	for ind, item := range rule.V1CanUse {
		canUse[ind] = speca.NewReferable(
			item,
			resolve(fmt.Sprintf("$.deps.%s.canUse[%d]", name, ind)),
		)
	}
	rule.internalCanUse = canUse

	// --
	mayDependOn := make([]speca.Referable[string], len(rule.V1MayDependOn))
	for ind, item := range rule.V1MayDependOn {
		mayDependOn[ind] = speca.NewReferable(
			item,
			resolve(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, ind)),
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

func (a archV1InternalCommonComponents) List() []speca.Referable[string] {
	return a.data
}

func (a archV1InternalCommonVendors) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalCommonVendors) List() []speca.Referable[string] {
	return a.data
}

func (a archV1InternalExcludeFilesRegExp) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalExcludeFilesRegExp) List() []speca.Referable[string] {
	return a.data
}

func (a archV1InternalExclude) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalExclude) List() []speca.Referable[string] {
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
