package arch

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type (
	// ComponentName is abstraction useful for mapping real packages to one Component.
	ComponentName = string

	// VendorName is abstraction useful for mapping real vendor packages to one Vendor.
	VendorName = string

	Document interface {
		Reference() common.Reference

		// FilePath is Abs arch file path
		FilePath() common.Referable[string]

		// Version of spec (scheme of document)
		Version() common.Referable[int]

		// WorkingDirectory relative to root, prepend this to all path's from spec
		WorkingDirectory() common.Referable[string]

		// Options is global spec options
		Options() Options

		// ExcludedDirectories from analyze, each contain relative directory name
		ExcludedDirectories() ExcludedDirectories

		// ExcludedFilesRegExp from analyze, each project file will be matched with this regexp rules
		ExcludedFilesRegExp() ExcludedFilesRegExp

		// Vendors (map)
		Vendors() Vendors

		// CommonVendors is list of Vendors that can be imported to any project package
		CommonVendors() CommonVendors

		// Components (map)
		Components() Components

		// CommonComponents is List of Components that can be imported to any project package
		CommonComponents() CommonComponents

		// Dependencies map between Components and DependencyRule`s
		Dependencies() Dependencies
	}

	Options interface {
		Reference() common.Reference

		// IsDependOnAnyVendor allows all project code depend on any third party vendor lib
		// analyze will not check imports with not local namespace's
		IsDependOnAnyVendor() common.Referable[bool]

		// DeepScan turn on usage of advanced AST linter
		// this is default behavior since v3+ configs
		DeepScan() common.Referable[bool]
	}

	ExcludedDirectories interface {
		Reference() common.Reference

		// List of directories
		// examples:
		// 	- internal/test
		//	- vendor
		//	- .idea
		List() []common.Referable[string]
	}

	ExcludedFilesRegExp interface {
		Reference() common.Reference

		// List of regexp's
		// examples:
		// 	- "^.*_test\\.go$"
		List() []common.Referable[string]
	}

	Vendors interface {
		Reference() common.Reference

		// Map describe Vendor packages properties
		Map() map[VendorName]Vendor
	}

	Vendor interface {
		Reference() common.Reference

		// ImportPaths is list of full import vendor qualified path
		// example:
		// 	- golang.org/x/mod/modfile
		// 	- example.com/*/libs/**
		ImportPaths() []common.Referable[models.Glob]
	}

	CommonVendors interface {
		Reference() common.Reference

		// List of Vendors that can be imported to any project package
		List() []common.Referable[string]
	}

	Components interface {
		Reference() common.Reference

		// Map with Component packages properties
		Map() map[ComponentName]Component
	}

	Component interface {
		Reference() common.Reference

		// RelativePaths can contain glob's
		// example:
		// 	- internal/service/*/models/**
		// 	- /
		// 	- tests/**
		RelativePaths() []common.Referable[models.Glob]
	}

	CommonComponents interface {
		Reference() common.Reference

		// List of Components that can be imported to any project package
		List() []common.Referable[string]
	}

	Dependencies interface {
		Reference() common.Reference

		// Map with Dependencies between Components and DependencyRule`s
		Map() map[ComponentName]DependencyRule
	}

	DependencyRule interface {
		Reference() common.Reference

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
