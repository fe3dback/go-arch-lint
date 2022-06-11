package spec

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ArchV2Document struct {
		reference                  models.Reference
		internalVersion            speca.Referable[int]
		internalWorkingDir         speca.Referable[string]
		internalVendors            archV2InternalVendors
		internalExclude            archV2InternalExclude
		internalExcludeFilesRegExp archV2InternalExcludeFilesRegExp
		internalComponents         archV2InternalComponents
		internalDependencies       archV2InternalDependencies
		internalCommonComponents   archV2InternalCommonComponents
		internalCommonVendors      archV2InternalCommonVendors

		V2Version            int                                    `yaml:"version" json:"version"`
		V2WorkDir            string                                 `yaml:"workdir" json:"workdir"`
		V2Allow              ArchV2Allow                            `yaml:"allow" json:"allow"`
		V2Vendors            map[arch.VendorName]ArchV2Vendor       `yaml:"vendors" json:"vendors"`
		V2Exclude            []string                               `yaml:"exclude" json:"exclude"`
		V2ExcludeFilesRegExp []string                               `yaml:"excludeFiles" json:"excludeFiles"`
		V2Components         map[arch.ComponentName]ArchV2Component `yaml:"components" json:"components"`
		V2Dependencies       map[arch.ComponentName]ArchV2Rules     `yaml:"deps" json:"deps"`
		V2CommonComponents   []string                               `yaml:"commonComponents" json:"commonComponents"`
		V2CommonVendors      []string                               `yaml:"commonVendors" json:"commonVendors"`
	}

	ArchV2Allow struct {
		reference              models.Reference
		internalDepOnAnyVendor speca.Referable[bool]

		V2DepOnAnyVendor bool `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
	}

	ArchV2Vendor struct {
		reference           models.Reference
		internalImportPaths []speca.Referable[models.Glob]

		V2ImportPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV2Component struct {
		reference          models.Reference
		internalLocalPaths []speca.Referable[models.Glob]

		V2LocalPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV2Rules struct {
		reference              models.Reference
		internalMayDependOn    []speca.Referable[string]
		internalCanUse         []speca.Referable[string]
		internalAnyProjectDeps speca.Referable[bool]
		internalAnyVendorDeps  speca.Referable[bool]

		V2MayDependOn    []string `yaml:"mayDependOn" json:"mayDependOn"`
		V2CanUse         []string `yaml:"canUse" json:"canUse"`
		V2AnyProjectDeps bool     `yaml:"anyProjectDeps" json:"anyProjectDeps"`
		V2AnyVendorDeps  bool     `yaml:"anyVendorDeps" json:"anyVendorDeps"`
	}
)

type (
	archV2InternalVendors struct {
		reference models.Reference
		data      map[arch.VendorName]ArchV2Vendor
	}

	archV2InternalComponents struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV2Component
	}

	archV2InternalExclude struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV2InternalExcludeFilesRegExp struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV2InternalCommonVendors struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV2InternalCommonComponents struct {
		reference models.Reference
		data      []speca.Referable[string]
	}

	archV2InternalDependencies struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV2Rules
	}
)

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (doc ArchV2Document) Reference() models.Reference {
	return doc.reference
}

func (doc ArchV2Document) Version() speca.Referable[int] {
	return doc.internalVersion
}

func (doc ArchV2Document) WorkingDirectory() speca.Referable[string] {
	return doc.internalWorkingDir
}

func (doc ArchV2Document) Options() arch.Options {
	return doc.V2Allow
}

func (doc ArchV2Document) ExcludedDirectories() arch.ExcludedDirectories {
	return doc.internalExclude
}

func (doc ArchV2Document) ExcludedFilesRegExp() arch.ExcludedFilesRegExp {
	return doc.internalExcludeFilesRegExp
}

func (doc ArchV2Document) Vendors() arch.Vendors {
	return doc.internalVendors
}

func (doc ArchV2Document) Components() arch.Components {
	return doc.internalComponents
}

func (doc ArchV2Document) CommonComponents() arch.CommonComponents {
	return doc.internalCommonComponents
}

func (doc ArchV2Document) CommonVendors() arch.CommonVendors {
	return doc.internalCommonVendors
}

func (doc ArchV2Document) Dependencies() arch.Dependencies {
	return doc.internalDependencies
}

func (doc ArchV2Document) applyReferences(resolver YAMLSourceCodeReferenceResolver) ArchV2Document {
	doc.reference = resolver.Resolve("$.version")

	// Version
	doc.internalVersion = speca.NewReferable(
		doc.V2Version,
		resolver.Resolve("$.version"),
	)

	// Working Directory
	actualWorkDirectory := "./" // fallback from version 1
	if doc.V2WorkDir != "" {
		actualWorkDirectory = doc.V2WorkDir
	}

	doc.internalWorkingDir = speca.NewReferable(
		actualWorkDirectory,
		resolver.Resolve("$.workdir"),
	)

	// Allow
	doc.V2Allow = doc.V2Allow.applyReferences(resolver)

	// Vendors
	vendors := make(map[string]ArchV2Vendor)
	for name, vendor := range doc.V2Vendors {
		vendors[name] = vendor.applyReferences(name, resolver)
	}
	doc.internalVendors = archV2InternalVendors{
		reference: resolver.Resolve("$.vendors"),
		data:      vendors,
	}

	// Exclude
	excludedDirectories := make([]speca.Referable[string], len(doc.V2Exclude))
	for ind, item := range doc.V2Exclude {
		excludedDirectories[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.exclude[%d]", ind)),
		)
	}

	doc.internalExclude = archV2InternalExclude{
		reference: resolver.Resolve("$.exclude"),
		data:      excludedDirectories,
	}

	// ExcludeFilesRegExp
	excludedFiles := make([]speca.Referable[string], len(doc.V2ExcludeFilesRegExp))
	for ind, item := range doc.V2ExcludeFilesRegExp {
		excludedFiles[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.excludeFiles[%d]", ind)),
		)
	}

	doc.internalExcludeFilesRegExp = archV2InternalExcludeFilesRegExp{
		reference: resolver.Resolve("$.excludeFiles"),
		data:      excludedFiles,
	}

	// Components
	components := make(map[string]ArchV2Component)
	for name, component := range doc.V2Components {
		components[name] = component.applyReferences(name, doc.internalWorkingDir.Value(), resolver)
	}
	doc.internalComponents = archV2InternalComponents{
		reference: resolver.Resolve("$.components"),
		data:      components,
	}

	// Dependencies
	dependencies := make(map[string]ArchV2Rules)
	for name, rules := range doc.V2Dependencies {
		dependencies[name] = rules.applyReferences(name, resolver)
	}
	doc.internalDependencies = archV2InternalDependencies{
		reference: resolver.Resolve("$.deps"),
		data:      dependencies,
	}

	// CommonComponents
	commonComponents := make([]speca.Referable[string], len(doc.V2CommonComponents))
	for ind, item := range doc.V2CommonComponents {
		commonComponents[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.commonComponents[%d]", ind)),
		)
	}
	doc.internalCommonComponents = archV2InternalCommonComponents{
		reference: resolver.Resolve("$.commonComponents"),
		data:      commonComponents,
	}

	// CommonVendors
	commonVendors := make([]speca.Referable[string], len(doc.V2CommonVendors))
	for ind, item := range doc.V2CommonVendors {
		commonVendors[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.commonVendors[%d]", ind)),
		)
	}
	doc.internalCommonVendors = archV2InternalCommonVendors{
		reference: resolver.Resolve("$.commonVendors"),
		data:      commonVendors,
	}

	return doc
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (opt ArchV2Allow) Reference() models.Reference {
	return opt.reference
}

func (opt ArchV2Allow) IsDependOnAnyVendor() speca.Referable[bool] {
	return opt.internalDepOnAnyVendor
}

func (opt ArchV2Allow) applyReferences(resolver YAMLSourceCodeReferenceResolver) ArchV2Allow {
	opt.reference = resolver.Resolve("$.allow")

	opt.internalDepOnAnyVendor = speca.NewReferable(
		opt.V2DepOnAnyVendor,
		resolver.Resolve("$.allow.depOnAnyVendor"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV2Vendor) Reference() models.Reference {
	return v.reference
}

func (v ArchV2Vendor) ImportPaths() []speca.Referable[models.Glob] {
	return v.internalImportPaths
}

func (v ArchV2Vendor) applyReferences(name arch.VendorName, resolver YAMLSourceCodeReferenceResolver) ArchV2Vendor {
	v.reference = resolver.Resolve(fmt.Sprintf("$.vendors.%s", name))

	for ind, importPath := range v.V2ImportPaths.list {
		yamlPath := fmt.Sprintf("$.vendors.%s.in", name)
		if v.V2ImportPaths.definedAsList {
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

func (c ArchV2Component) Reference() models.Reference {
	return c.reference
}

func (c ArchV2Component) RelativePaths() []speca.Referable[models.Glob] {
	return c.internalLocalPaths
}

func (c ArchV2Component) applyReferences(
	name arch.ComponentName,
	workDirectory string,
	resolver YAMLSourceCodeReferenceResolver,
) ArchV2Component {
	c.reference = resolver.Resolve(fmt.Sprintf("$.components.%s", name))

	for ind, importPath := range c.V2LocalPaths.list {
		yamlPath := fmt.Sprintf("$.components.%s.in", name)
		if c.V2LocalPaths.definedAsList {
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

func (rule ArchV2Rules) Reference() models.Reference {
	return rule.reference
}

func (rule ArchV2Rules) MayDependOn() []speca.Referable[string] {
	return rule.internalMayDependOn
}

func (rule ArchV2Rules) CanUse() []speca.Referable[string] {
	return rule.internalCanUse
}

func (rule ArchV2Rules) AnyProjectDeps() speca.Referable[bool] {
	return rule.internalAnyProjectDeps
}

func (rule ArchV2Rules) AnyVendorDeps() speca.Referable[bool] {
	return rule.internalAnyVendorDeps
}

func (rule ArchV2Rules) applyReferences(name arch.ComponentName, resolver YAMLSourceCodeReferenceResolver) ArchV2Rules {
	rule.reference = resolver.Resolve(fmt.Sprintf("$.deps.%s", name))

	// --
	rule.internalAnyVendorDeps = speca.NewReferable(
		rule.V2AnyVendorDeps,
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyVendorDeps", name)),
	)

	// --
	rule.internalAnyProjectDeps = speca.NewReferable(
		rule.V2AnyProjectDeps,
		resolver.Resolve(fmt.Sprintf("$.deps.%s.anyProjectDeps", name)),
	)

	// --
	canUse := make([]speca.Referable[string], len(rule.V2CanUse))
	for ind, item := range rule.V2CanUse {
		canUse[ind] = speca.NewReferable(
			item,
			resolver.Resolve(fmt.Sprintf("$.deps.%s.canUse[%d]", name, ind)),
		)
	}
	rule.internalCanUse = canUse

	// --
	mayDependOn := make([]speca.Referable[string], len(rule.V2MayDependOn))
	for ind, item := range rule.V2MayDependOn {
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

func (a archV2InternalDependencies) Reference() models.Reference {
	return a.reference
}

func (a archV2InternalDependencies) Map() map[arch.ComponentName]arch.DependencyRule {
	res := make(map[arch.ComponentName]arch.DependencyRule)
	for name, rules := range a.data {
		res[name] = rules
	}
	return res
}

func (a archV2InternalCommonComponents) Reference() models.Reference {
	return a.reference
}

func (a archV2InternalCommonComponents) List() []speca.Referable[string] {
	return a.data
}

func (a archV2InternalCommonVendors) Reference() models.Reference {
	return a.reference
}

func (a archV2InternalCommonVendors) List() []speca.Referable[string] {
	return a.data
}

func (a archV2InternalExcludeFilesRegExp) Reference() models.Reference {
	return a.reference
}

func (a archV2InternalExcludeFilesRegExp) List() []speca.Referable[string] {
	return a.data
}

func (a archV2InternalExclude) Reference() models.Reference {
	return a.reference
}

func (a archV2InternalExclude) List() []speca.Referable[string] {
	return a.data
}

func (a archV2InternalComponents) Reference() models.Reference {
	return a.reference
}

func (a archV2InternalComponents) Map() map[arch.ComponentName]arch.Component {
	res := make(map[arch.ComponentName]arch.Component)
	for name, component := range a.data {
		res[name] = component
	}
	return res
}

func (a archV2InternalVendors) Reference() models.Reference {
	return a.reference
}

func (a archV2InternalVendors) Map() map[arch.VendorName]arch.Vendor {
	res := make(map[arch.VendorName]arch.Vendor)
	for name, vendor := range a.data {
		res[name] = vendor
	}
	return res
}
