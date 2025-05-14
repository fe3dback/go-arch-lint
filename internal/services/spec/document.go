package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

// spec layout:
//   decoder   - decode different config file formats (v1 ... latest) into one single Document interface
//   validator - will validate Document interface (check integrity between config fields)
//   assembler - will assemble arch.Spec from validated Document interface
//
// all other operations (business logic) code will use arch.Spec object for reading config values.

type (
	// ComponentName is abstraction useful for mapping real packages to one Component.
	ComponentName = string

	// VendorName is abstraction useful for mapping real vendor packages to one Vendor.
	VendorName = string

	Vendors      = map[VendorName]common.Referable[Vendor]
	Components   = map[ComponentName]common.Referable[Component]
	Dependencies = map[ComponentName]common.Referable[DependencyRule]

	Document interface {
		// Version of spec (scheme of document)
		Version() common.Referable[int]

		// WorkingDirectory relative to root, prepend this to all path's from spec
		WorkingDirectory() common.Referable[string]

		// Options is global spec options
		Options() Options

		// ExcludedDirectories from analyze, each contain relative directory name
		// List of directories
		// examples:
		// 	- internal/test
		//	- vendor
		//	- .idea
		ExcludedDirectories() []common.Referable[string]

		// ExcludedFilesRegExp from analyze, each project file will be matched with this regexp rules
		// List of regexp's
		// examples:
		// 	- "^.*_test\\.go$"
		ExcludedFilesRegExp() []common.Referable[string]

		// Vendors (map)
		Vendors() Vendors

		// CommonVendors is list of Vendors that can be imported to any project package
		CommonVendors() []common.Referable[string]

		// Components (map)
		Components() Components

		// CommonComponents is List of Components that can be imported to any project package
		CommonComponents() []common.Referable[string]

		// Dependencies map between Components and DependencyRule`s
		Dependencies() Dependencies
	}

	Options interface {
		// IsDependOnAnyVendor allows all project code depend on any third party vendor lib
		// analyze will not check imports with not local namespace's
		IsDependOnAnyVendor() common.Referable[bool]

		// DeepScan turn on usage of advanced AST linter
		// this is default behavior since v3+ configs
		DeepScan() common.Referable[bool]

		// IgnoreNotFoundComponents skips components that are not found by their glob
		// disabled by default
		IgnoreNotFoundComponents() common.Referable[bool]
	}

	Vendor interface {
		// ImportPaths is list of full import vendor qualified path
		// example:
		// 	- golang.org/x/mod/modfile
		// 	- example.com/*/libs/**
		ImportPaths() []models.Glob
	}

	Component interface {
		// RelativePaths can contain glob's
		// example:
		// 	- internal/service/*/models/**
		// 	- /
		// 	- tests/**
		RelativePaths() []models.Glob
	}

	DependencyRule interface {
		// MayDependOn is list of Component names, that can be imported to described component
		MayDependOn() []common.Referable[string]

		// CanUse is list of Vendor names, that can be imported to described component
		CanUse() []common.Referable[string]

		// AnyProjectDeps allow component to import any other local namespace packages
		AnyProjectDeps() common.Referable[bool]

		// AnyVendorDeps allow component to import any other vendor namespace packages
		AnyVendorDeps() common.Referable[bool]

		// DeepScan overrides deepScan global option
		DeepScan() common.Referable[bool]
	}
)
