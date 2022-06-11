package arch

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	// abstract ComponentName useful for mapping real packages to one Component.
	ComponentName = string

	// abstract VendorName useful for mapping real vendor packages to one Vendor.
	VendorName = string

	Document interface {
		Reference() models.Reference

		// Spec Version
		Version() speca.Referable[int]

		// Spec relative WorkingDirectory to root, prepend this to all path's from spec
		WorkingDirectory() speca.Referable[string]

		// Global spec Options
		Options() Options

		// ExcludedDirectories from analyze, each contain relative directory name
		ExcludedDirectories() ExcludedDirectories

		// ExcludedFilesRegExp from analyze, each project file will by matched with this regexp rules
		ExcludedFilesRegExp() ExcludedFilesRegExp

		// Map of Vendors
		Vendors() Vendors

		// List of Vendors that can by imported to any project package
		CommonVendors() CommonVendors

		// Map of Components
		Components() Components

		// List of Components that can by imported to any project package
		CommonComponents() CommonComponents

		// Dependencies map between Components and DependencyRule`s
		Dependencies() Dependencies
	}

	Options interface {
		Reference() models.Reference

		// allow all project code depend on any third party vendor lib
		// analyze will not check imports with not local namespace's
		IsDependOnAnyVendor() speca.Referable[bool]
	}

	ExcludedDirectories interface {
		Reference() models.Reference

		// list of directories
		// examples:
		// 	- internal/test
		//	- vendor
		//	- .idea
		List() []speca.Referable[string]
	}

	ExcludedFilesRegExp interface {
		Reference() models.Reference

		// list of regexp's
		// examples:
		// 	- "^.*_test\\.go$"
		List() []speca.Referable[string]
	}

	Vendors interface {
		Reference() models.Reference

		// describe Vendor packages properties
		Map() map[VendorName]Vendor
	}

	Vendor interface {
		Reference() models.Reference

		// Full import vendor qualified path
		// example:
		// 	- golang.org/x/mod/modfile
		// 	- example.com/*/libs/**
		ImportPaths() []speca.Referable[models.Glob]
	}

	CommonVendors interface {
		Reference() models.Reference

		// List of Vendors that can by imported to any project package
		List() []speca.Referable[string]
	}

	Components interface {
		Reference() models.Reference

		// describe Component packages properties
		Map() map[ComponentName]Component
	}

	Component interface {
		Reference() models.Reference

		// Relative package path, can contain glob's
		// example:
		// 	- internal/service/*/models/**
		// 	- /
		// 	- tests/**
		RelativePaths() []speca.Referable[models.Glob]
	}

	CommonComponents interface {
		Reference() models.Reference

		// List of Components that can by imported to any project package
		List() []speca.Referable[string]
	}

	Dependencies interface {
		Reference() models.Reference

		// Dependencies map between Components and DependencyRule`s
		Map() map[ComponentName]DependencyRule
	}

	DependencyRule interface {
		Reference() models.Reference

		// List of Component names, that can by imported to described component
		MayDependOn() []speca.Referable[string]

		// List of Vendor names, that can by imported to described component
		CanUse() []speca.Referable[string]

		// described component can import any other local namespace packages
		AnyProjectDeps() speca.Referable[bool]

		// described component can import any other vendor namespace packages
		AnyVendorDeps() speca.Referable[bool]
	}
)
