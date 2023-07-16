package spec

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type (
	ArchV3Document struct {
		filePath common.Referable[string]

		reference                  common.Reference
		internalVersion            common.Referable[int]
		internalWorkingDir         common.Referable[string]
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
		reference              common.Reference
		internalDepOnAnyVendor common.Referable[bool]
		internalDeepScan       common.Referable[bool]

		V3DepOnAnyVendor bool  `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
		V3DeepScan       *bool `yaml:"deepScan" json:"deepScan"`
	}

	ArchV3Vendor struct {
		reference           common.Reference
		internalImportPaths []common.Referable[models.Glob]

		V3ImportPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV3Component struct {
		reference          common.Reference
		internalLocalPaths []common.Referable[models.Glob]

		V3LocalPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV3Rules struct {
		reference              common.Reference
		internalMayDependOn    []common.Referable[string]
		internalCanUse         []common.Referable[string]
		internalAnyProjectDeps common.Referable[bool]
		internalAnyVendorDeps  common.Referable[bool]
		internalDeepScan       common.Referable[bool]

		V3MayDependOn    []string `yaml:"mayDependOn" json:"mayDependOn"`
		V3CanUse         []string `yaml:"canUse" json:"canUse"`
		V3AnyProjectDeps bool     `yaml:"anyProjectDeps" json:"anyProjectDeps"`
		V3AnyVendorDeps  bool     `yaml:"anyVendorDeps" json:"anyVendorDeps"`
		V3DeepScan       *bool    `yaml:"deepScan" json:"deepScan"`
	}
)

type (
	archV3InternalVendors struct {
		reference common.Reference
		data      map[arch.VendorName]ArchV3Vendor
	}

	archV3InternalComponents struct {
		reference common.Reference
		data      map[arch.ComponentName]ArchV3Component
	}

	archV3InternalExclude struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV3InternalExcludeFilesRegExp struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV3InternalCommonVendors struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV3InternalCommonComponents struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV3InternalDependencies struct {
		reference common.Reference
		data      map[arch.ComponentName]ArchV3Rules
	}
)

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (doc ArchV3Document) FilePath() common.Referable[string] {
	return doc.filePath
}

func (doc ArchV3Document) Reference() common.Reference {
	return doc.reference
}

func (doc ArchV3Document) Version() common.Referable[int] {
	return doc.internalVersion
}

func (doc ArchV3Document) WorkingDirectory() common.Referable[string] {
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

func (doc ArchV3Document) applyReferences(resolve yamlDocumentPathResolver) ArchV3Document {
	doc.reference = resolve("$.version")

	// Version
	doc.internalVersion = common.NewReferable(
		doc.V3Version,
		resolve("$.version"),
	)

	// Working Directory
	actualWorkDirectory := "./" // fallback from version 1
	if doc.V3WorkDir != "" {
		actualWorkDirectory = doc.V3WorkDir
	}

	doc.internalWorkingDir = common.NewReferable(
		actualWorkDirectory,
		resolve("$.workdir"),
	)

	// Allow
	doc.V3Allow = doc.V3Allow.applyReferences(resolve)

	// Vendors
	vendors := make(map[string]ArchV3Vendor)
	for name, vendor := range doc.V3Vendors {
		vendors[name] = vendor.applyReferences(name, resolve)
	}
	doc.internalVendors = archV3InternalVendors{
		reference: resolve("$.vendors"),
		data:      vendors,
	}

	// Exclude
	excludedDirectories := make([]common.Referable[string], len(doc.V3Exclude))
	for ind, item := range doc.V3Exclude {
		excludedDirectories[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.exclude[%d]", ind)),
		)
	}

	doc.internalExclude = archV3InternalExclude{
		reference: resolve("$.exclude"),
		data:      excludedDirectories,
	}

	// ExcludeFilesRegExp
	excludedFiles := make([]common.Referable[string], len(doc.V3ExcludeFilesRegExp))
	for ind, item := range doc.V3ExcludeFilesRegExp {
		excludedFiles[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.excludeFiles[%d]", ind)),
		)
	}

	doc.internalExcludeFilesRegExp = archV3InternalExcludeFilesRegExp{
		reference: resolve("$.excludeFiles"),
		data:      excludedFiles,
	}

	// Components
	components := make(map[string]ArchV3Component)
	for name, component := range doc.V3Components {
		components[name] = component.applyReferences(name, doc.internalWorkingDir.Value, resolve)
	}
	doc.internalComponents = archV3InternalComponents{
		reference: resolve("$.components"),
		data:      components,
	}

	// Dependencies
	dependencies := make(map[string]ArchV3Rules)
	for name, rules := range doc.V3Dependencies {
		dependencies[name] = rules.applyReferences(name, doc.V3Allow, resolve)
	}
	doc.internalDependencies = archV3InternalDependencies{
		reference: resolve("$.deps"),
		data:      dependencies,
	}

	// CommonComponents
	commonComponents := make([]common.Referable[string], len(doc.V3CommonComponents))
	for ind, item := range doc.V3CommonComponents {
		commonComponents[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.commonComponents[%d]", ind)),
		)
	}
	doc.internalCommonComponents = archV3InternalCommonComponents{
		reference: resolve("$.commonComponents"),
		data:      commonComponents,
	}

	// CommonVendors
	commonVendors := make([]common.Referable[string], len(doc.V3CommonVendors))
	for ind, item := range doc.V3CommonVendors {
		commonVendors[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.commonVendors[%d]", ind)),
		)
	}
	doc.internalCommonVendors = archV3InternalCommonVendors{
		reference: resolve("$.commonVendors"),
		data:      commonVendors,
	}

	return doc
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (opt ArchV3Allow) Reference() common.Reference {
	return opt.reference
}

func (opt ArchV3Allow) IsDependOnAnyVendor() common.Referable[bool] {
	return opt.internalDepOnAnyVendor
}

func (opt ArchV3Allow) DeepScan() common.Referable[bool] {
	return opt.internalDeepScan
}

func (opt ArchV3Allow) applyReferences(resolve yamlDocumentPathResolver) ArchV3Allow {
	opt.reference = resolve("$.allow")

	opt.internalDepOnAnyVendor = common.NewReferable(
		opt.V3DepOnAnyVendor,
		resolve("$.allow.depOnAnyVendor"),
	)

	deepScan := false
	if opt.V3DeepScan == nil {
		// be default it`s on from V3+
		deepScan = true
	} else {
		deepScan = *opt.V3DeepScan
	}

	opt.internalDeepScan = common.NewReferable(
		deepScan,
		resolve("$.allow.deepScan"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV3Vendor) Reference() common.Reference {
	return v.reference
}

func (v ArchV3Vendor) ImportPaths() []common.Referable[models.Glob] {
	return v.internalImportPaths
}

func (v ArchV3Vendor) applyReferences(name arch.VendorName, resolve yamlDocumentPathResolver) ArchV3Vendor {
	v.reference = resolve(fmt.Sprintf("$.vendors.%s", name))

	for ind, importPath := range v.V3ImportPaths.list {
		yamlPath := fmt.Sprintf("$.vendors.%s.in", name)
		if v.V3ImportPaths.definedAsList {
			yamlPath = fmt.Sprintf("%s[%d]", yamlPath, ind)
		}

		v.internalImportPaths = append(v.internalImportPaths, common.NewReferable(
			models.Glob(importPath),
			resolve(yamlPath),
		))
	}

	return v
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (c ArchV3Component) Reference() common.Reference {
	return c.reference
}

func (c ArchV3Component) RelativePaths() []common.Referable[models.Glob] {
	return c.internalLocalPaths
}

func (c ArchV3Component) applyReferences(
	name arch.ComponentName,
	workDirectory string,
	resolve yamlDocumentPathResolver,
) ArchV3Component {
	c.reference = resolve(fmt.Sprintf("$.components.%s", name))

	for ind, importPath := range c.V3LocalPaths.list {
		yamlPath := fmt.Sprintf("$.components.%s.in", name)
		if c.V3LocalPaths.definedAsList {
			yamlPath = fmt.Sprintf("%s[%d]", yamlPath, ind)
		}

		c.internalLocalPaths = append(c.internalLocalPaths, common.NewReferable(
			models.Glob(
				path.Clean(fmt.Sprintf("%s/%s",
					workDirectory,
					importPath,
				)),
			),
			resolve(yamlPath),
		))
	}

	return c
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (rule ArchV3Rules) Reference() common.Reference {
	return rule.reference
}

func (rule ArchV3Rules) MayDependOn() []common.Referable[string] {
	return rule.internalMayDependOn
}

func (rule ArchV3Rules) CanUse() []common.Referable[string] {
	return rule.internalCanUse
}

func (rule ArchV3Rules) AnyProjectDeps() common.Referable[bool] {
	return rule.internalAnyProjectDeps
}

func (rule ArchV3Rules) AnyVendorDeps() common.Referable[bool] {
	return rule.internalAnyVendorDeps
}

func (rule ArchV3Rules) DeepScan() common.Referable[bool] {
	return rule.internalDeepScan
}

func (rule ArchV3Rules) applyReferences(name arch.ComponentName, globalOptions ArchV3Allow, resolve yamlDocumentPathResolver) ArchV3Rules {
	rule.reference = resolve(fmt.Sprintf("$.deps.%s", name))

	// --
	if rule.V3DeepScan == nil {
		rule.internalDeepScan = globalOptions.internalDeepScan
	} else {
		// override deepScan for this component
		rule.internalDeepScan = common.NewReferable(
			*rule.V3DeepScan,
			resolve(fmt.Sprintf("$.deps.%s.deepScan", name)),
		)
	}

	// --
	rule.internalAnyVendorDeps = common.NewReferable(
		rule.V3AnyVendorDeps,
		resolve(fmt.Sprintf("$.deps.%s.anyVendorDeps", name)),
	)

	// --
	rule.internalAnyProjectDeps = common.NewReferable(
		rule.V3AnyProjectDeps,
		resolve(fmt.Sprintf("$.deps.%s.anyProjectDeps", name)),
	)

	// --
	canUse := make([]common.Referable[string], len(rule.V3CanUse))
	for ind, item := range rule.V3CanUse {
		canUse[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.deps.%s.canUse[%d]", name, ind)),
		)
	}
	rule.internalCanUse = canUse

	// --
	mayDependOn := make([]common.Referable[string], len(rule.V3MayDependOn))
	for ind, item := range rule.V3MayDependOn {
		mayDependOn[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, ind)),
		)
	}
	rule.internalMayDependOn = mayDependOn

	// --
	return rule
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (a archV3InternalDependencies) Reference() common.Reference {
	return a.reference
}

func (a archV3InternalDependencies) Map() map[arch.ComponentName]arch.DependencyRule {
	res := make(map[arch.ComponentName]arch.DependencyRule)
	for name, rules := range a.data {
		res[name] = rules
	}
	return res
}

func (a archV3InternalCommonComponents) Reference() common.Reference {
	return a.reference
}

func (a archV3InternalCommonComponents) List() []common.Referable[string] {
	return a.data
}

func (a archV3InternalCommonVendors) Reference() common.Reference {
	return a.reference
}

func (a archV3InternalCommonVendors) List() []common.Referable[string] {
	return a.data
}

func (a archV3InternalExcludeFilesRegExp) Reference() common.Reference {
	return a.reference
}

func (a archV3InternalExcludeFilesRegExp) List() []common.Referable[string] {
	return a.data
}

func (a archV3InternalExclude) Reference() common.Reference {
	return a.reference
}

func (a archV3InternalExclude) List() []common.Referable[string] {
	return a.data
}

func (a archV3InternalComponents) Reference() common.Reference {
	return a.reference
}

func (a archV3InternalComponents) Map() map[arch.ComponentName]arch.Component {
	res := make(map[arch.ComponentName]arch.Component)
	for name, component := range a.data {
		res[name] = component
	}
	return res
}

func (a archV3InternalVendors) Reference() common.Reference {
	return a.reference
}

func (a archV3InternalVendors) Map() map[arch.VendorName]arch.Vendor {
	res := make(map[arch.VendorName]arch.Vendor)
	for name, vendor := range a.data {
		res[name] = vendor
	}
	return res
}
