package spec

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ArchV3Document struct {
		reference                  models.Reference
		internalVersion            speca.Referable[int]
		internalWorkingDir         speca.Referable[string]
		internalVendors            archV3InternalVendors
		internalExclude            archV3InternalExclude
		internalExcludeFilesRegExp archV3InternalExcludeFilesRegExp
		internalComponents         archV3InternalComponents
		internalDependencies       archV3InternalDependencies
		internalCommonComponents   archV3InternalCommonComponents
		internalCommonVendors      archV3InternalCommonVendors

		V3Version            int                                    `yaml:"version" json:"version"`
		V3WorkDir            string                                 `yaml:"workdir" json:"workdir"`
		V3Allow              ArchV3Allow                            `yaml:"allow" json:"allow"`
		V3Vendors            map[arch.VendorName]ArchV3Vendor       `yaml:"vendors" json:"vendors"`
		V3Exclude            []string                               `yaml:"exclude" json:"exclude"`
		V3ExcludeFilesRegExp []string                               `yaml:"excludeFiles" json:"excludeFiles"`
		V3Components         map[arch.ComponentName]ArchV3Component `yaml:"components" json:"components"`
		V3Dependencies       map[arch.ComponentName]ArchV3Rules     `yaml:"deps" json:"deps"`
		V3CommonComponents   []string                               `yaml:"commonComponents" json:"commonComponents"`
		V3CommonVendors      []string                               `yaml:"commonVendors" json:"commonVendors"`
	}

	ArchV3Allow struct {
		reference              models.Reference
		internalDepOnAnyVendor speca.Referable[bool]
		internalDeepScan       speca.Referable[bool]

		V3DepOnAnyVendor bool  `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
		V3DeepScan       *bool `yaml:"deepScan" json:"deepScan"`
	}

	ArchV3Vendor struct {
		reference           models.Reference
		internalImportPaths []speca.Referable[models.Glob]

		V3ImportPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV3Component struct {
		reference          models.Reference
		internalLocalPaths []speca.Referable[models.Glob]

		V3LocalPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV3Rules struct {
		reference              models.Reference
		internalMayDependOn    []speca.Referable[string]
		internalCanUse         []speca.Referable[string]
		internalAnyProjectDeps speca.Referable[bool]
		internalAnyVendorDeps  speca.Referable[bool]
		internalDeepScan       speca.Referable[bool]

		V3MayDependOn    []string `yaml:"mayDependOn" json:"mayDependOn"`
		V3CanUse         []string `yaml:"canUse" json:"canUse"`
		V3AnyProjectDeps bool     `yaml:"anyProjectDeps" json:"anyProjectDeps"`
		V3AnyVendorDeps  bool     `yaml:"anyVendorDeps" json:"anyVendorDeps"`
		V3DeepScan       *bool    `yaml:"deepScan" json:"deepScan"`
	}
)

type (
	archV3InternalVendors struct {
		reference models.Reference
		data      map[arch.VendorName]ArchV3Vendor
	}

	archV3InternalComponents struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV3Component
	}

	archV3InternalExclude struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV3InternalExcludeFilesRegExp struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV3InternalCommonVendors struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV3InternalCommonComponents struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV3InternalDependencies struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV3Rules
	}
)

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (doc ArchV3Document) Reference() models.Reference {
	return doc.reference
}

func (doc ArchV3Document) Version() speca.Referable[int] {
	return doc.internalVersion
}

func (doc ArchV3Document) WorkingDirectory() speca.Referable[string] {
	return doc.internalWorkingDir
}

func (doc ArchV3Document) Options() arch.Options {
	return doc.V3Allow
}

func (doc ArchV3Document) ExcludedDirectories() arch.ExcludedDirectories {
	return doc.internalExclude
}

func (doc ArchV3Document) ExcludedFilesRegExp() arch.ExcludedFilesRegExp {
	return doc.internalExcludeFilesRegExp
}

func (doc ArchV3Document) Vendors() arch.Vendors {
	return doc.internalVendors
}

func (doc ArchV3Document) Components() arch.Components {
	return doc.internalComponents
}

func (doc ArchV3Document) CommonComponents() arch.CommonComponents {
	return doc.internalCommonComponents
}

func (doc ArchV3Document) CommonVendors() arch.CommonVendors {
	return doc.internalCommonVendors
}

func (doc ArchV3Document) Dependencies() arch.Dependencies {
	return doc.internalDependencies
}

func (doc ArchV3Document) applyReferences(resolver YAMLSourceCodeReferenceResolver) ArchV3Document {
	doc.reference = resolver.Resolve("$.version")

	// Version
	doc.internalVersion = speca.NewReferable(
		doc.V3Version,
		resolver.Resolve("$.version"),
	)

	// Working Directory
	actualWorkDirectory := "./" // fallback from version 1
	if doc.V3WorkDir != "" {
		actualWorkDirectory = doc.V3WorkDir
	}

	doc.internalWorkingDir = speca.NewReferable(
		actualWorkDirectory,
		resolver.Resolve("$.workdir"),
	)

	// Allow
	doc.V3Allow = doc.V3Allow.applyReferences(resolver)

	// Vendors
	vendors := make(map[string]ArchV3Vendor)
	for name, vendor := range doc.V3Vendors {
		vendors[name] = vendor.applyReferences(name, resolver)
	}
	doc.internalVendors = archV3InternalVendors{
		reference: resolver.Resolve("$.vendors"),
		data:      vendors,
	}

	// Exclude
	excludedDirectories := make([]speca.Referable[string], len(doc.V3Exclude))
	for ind, item := range doc.V3Exclude {
		excludedDirectories[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.exclude[%d]", ind)),
		)
	}

	doc.internalExclude = archV3InternalExclude{
		reference: resolver.Resolve("$.exclude"),
		data:      excludedDirectories,
	}

	// ExcludeFilesRegExp
	excludedFiles := make([]speca.Referable[string], len(doc.V3ExcludeFilesRegExp))
	for ind, item := range doc.V3ExcludeFilesRegExp {
		excludedFiles[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.excludeFiles[%d]", ind)),
		)
	}

	doc.internalExcludeFilesRegExp = archV3InternalExcludeFilesRegExp{
		reference: resolver.Resolve("$.excludeFiles"),
		data:      excludedFiles,
	}

	// Components
	components := make(map[string]ArchV3Component)
	for name, component := range doc.V3Components {
		components[name] = component.applyReferences(name, doc.internalWorkingDir.Value(), resolver)
	}
	doc.internalComponents = archV3InternalComponents{
		reference: resolver.Resolve("$.components"),
		data:      components,
	}

	// Dependencies
	dependencies := make(map[string]ArchV3Rules)
	for name, rules := range doc.V3Dependencies {
		dependencies[name] = rules.applyReferences(name, doc.V3Allow, resolver)
	}
	doc.internalDependencies = archV3InternalDependencies{
		reference: resolver.Resolve("$.deps"),
		data:      dependencies,
	}

	// CommonComponents
	commonComponents := make([]speca.Referable[string], len(doc.V3CommonComponents))
	for ind, item := range doc.V3CommonComponents {
		commonComponents[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.commonComponents[%d]", ind)),
		)
	}
	doc.internalCommonComponents = archV3InternalCommonComponents{
		reference: resolver.Resolve("$.commonComponents"),
		data:      commonComponents,
	}

	// CommonVendors
	commonVendors := make([]speca.Referable[string], len(doc.V3CommonVendors))
	for ind, item := range doc.V3CommonVendors {
		commonVendors[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.commonVendors[%d]", ind)),
		)
	}
	doc.internalCommonVendors = archV3InternalCommonVendors{
		reference: resolver.Resolve("$.commonVendors"),
		data:      commonVendors,
	}

	return doc
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (opt ArchV3Allow) Reference() models.Reference {
	return opt.reference
}

func (opt ArchV3Allow) IsDependOnAnyVendor() speca.Referable[bool] {
	return opt.internalDepOnAnyVendor
}

func (opt ArchV3Allow) DeepScan() speca.Referable[bool] {
	return opt.internalDeepScan
}

func (opt ArchV3Allow) applyReferences(resolver YAMLSourceCodeReferenceResolver) ArchV3Allow {
	opt.reference = resolver.Resolve("$.allow")

	opt.internalDepOnAnyVendor = speca.NewReferable(
		opt.V3DepOnAnyVendor,
		resolver.Resolve("$.allow.depOnAnyVendor"),
	)

	deepScan := false
	if opt.V3DeepScan == nil {
		// be default it`s on from V3+
		deepScan = true
	} else {
		deepScan = *opt.V3DeepScan
	}

	opt.internalDeepScan = speca.NewReferable(
		deepScan,
		resolver.Resolve("$.allow.deepScan"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV3Vendor) Reference() models.Reference {
	return v.reference
}

func (v ArchV3Vendor) ImportPaths() []speca.Referable[models.Glob] {
	return v.internalImportPaths
}

func (v ArchV3Vendor) applyReferences(name arch.VendorName, resolver YAMLSourceCodeReferenceResolver) ArchV3Vendor {
	v.reference = resolver.Resolve(fmt.Sprintf("$.vendors.%s", name))

	for ind, importPath := range v.V3ImportPaths.list {
		yamlPath := fmt.Sprintf("$.vendors.%s.in", name)
		if v.V3ImportPaths.definedAsList {
			yamlPath = fmt.Sprintf("%s[%d]", yamlPath, ind)
		}

		v.internalImportPaths = append(v.internalImportPaths, speca.NewReferable(
			models.Glob(importPath),
			resolver.Resolve(yamlPath),
		))
	}

	return v
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (c ArchV3Component) Reference() models.Reference {
	return c.reference
}

func (c ArchV3Component) RelativePaths() []speca.Referable[models.Glob] {
	return c.internalLocalPaths
}

func (c ArchV3Component) applyReferences(
	name arch.ComponentName,
	workDirectory string,
	resolver YAMLSourceCodeReferenceResolver,
) ArchV3Component {
	c.reference = resolver.Resolve(fmt.Sprintf("$.components.%s", name))

	for ind, importPath := range c.V3LocalPaths.list {
		yamlPath := fmt.Sprintf("$.components.%s.in", name)
		if c.V3LocalPaths.definedAsList {
			yamlPath = fmt.Sprintf("%s[%d]", yamlPath, ind)
		}

		c.internalLocalPaths = append(c.internalLocalPaths, speca.NewReferable(
			models.Glob(
				path.Clean(fmt.Sprintf("%s/%s",
					workDirectory,
					importPath,
				)),
			),
			resolver.Resolve(yamlPath),
		))
	}

	return c
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (rule ArchV3Rules) Reference() models.Reference {
	return rule.reference
}

func (rule ArchV3Rules) MayDependOn() []speca.Referable[string] {
	return rule.internalMayDependOn
}

func (rule ArchV3Rules) CanUse() []speca.Referable[string] {
	return rule.internalCanUse
}

func (rule ArchV3Rules) AnyProjectDeps() speca.Referable[bool] {
	return rule.internalAnyProjectDeps
}

func (rule ArchV3Rules) AnyVendorDeps() speca.Referable[bool] {
	return rule.internalAnyVendorDeps
}

func (rule ArchV3Rules) DeepScan() speca.Referable[bool] {
	return rule.internalDeepScan
}

func (rule ArchV3Rules) applyReferences(name arch.ComponentName, globalOptions ArchV3Allow, resolver YAMLSourceCodeReferenceResolver) ArchV3Rules {
	rule.reference = resolver.Resolve(fmt.Sprintf("$.deps.%s", name))

	// --
	if rule.V3DeepScan == nil {
		rule.internalDeepScan = globalOptions.internalDeepScan
	} else {
		// override deepScan for this component
		rule.internalDeepScan = speca.NewReferable(
			*rule.V3DeepScan,
			resolver.Resolve(fmt.Sprintf("$.deps.%s.deepScan", name)),
		)
	}

	// --
	rule.internalAnyVendorDeps = speca.NewReferable(
		rule.V3AnyVendorDeps,
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyVendorDeps", name)),
	)

	// --
	rule.internalAnyProjectDeps = speca.NewReferable(
		rule.V3AnyProjectDeps,
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyProjectDeps", name)),
	)

	// --
	canUse := make([]speca.Referable[string], len(rule.V3CanUse))
	for ind, item := range rule.V3CanUse {
		canUse[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.deps.%s.canUse[%d]", name, ind)),
		)
	}
	rule.internalCanUse = canUse

	// --
	mayDependOn := make([]speca.Referable[string], len(rule.V3MayDependOn))
	for ind, item := range rule.V3MayDependOn {
		mayDependOn[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, ind)),
		)
	}
	rule.internalMayDependOn = mayDependOn

	// --
	return rule
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (a archV3InternalDependencies) Reference() models.Reference {
	return a.reference
}

func (a archV3InternalDependencies) Map() map[arch.ComponentName]arch.DependencyRule {
	res := make(map[arch.ComponentName]arch.DependencyRule)
	for name, rules := range a.data {
		res[name] = rules
	}
	return res
}

func (a archV3InternalCommonComponents) Reference() models.Reference {
	return a.reference
}

func (a archV3InternalCommonComponents) List() []speca.Referable[string] {
	return a.data
}

func (a archV3InternalCommonVendors) Reference() models.Reference {
	return a.reference
}

func (a archV3InternalCommonVendors) List() []speca.Referable[string] {
	return a.data
}

func (a archV3InternalExcludeFilesRegExp) Reference() models.Reference {
	return a.reference
}

func (a archV3InternalExcludeFilesRegExp) List() []speca.Referable[string] {
	return a.data
}

func (a archV3InternalExclude) Reference() models.Reference {
	return a.reference
}

func (a archV3InternalExclude) List() []speca.Referable[string] {
	return a.data
}

func (a archV3InternalComponents) Reference() models.Reference {
	return a.reference
}

func (a archV3InternalComponents) Map() map[arch.ComponentName]arch.Component {
	res := make(map[arch.ComponentName]arch.Component)
	for name, component := range a.data {
		res[name] = component
	}
	return res
}

func (a archV3InternalVendors) Reference() models.Reference {
	return a.reference
}

func (a archV3InternalVendors) Map() map[arch.VendorName]arch.Vendor {
	res := make(map[arch.VendorName]arch.Vendor)
	for name, vendor := range a.data {
		res[name] = vendor
	}
	return res
}
