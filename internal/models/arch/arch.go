package arch

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	// abstract ComponentName useful for mapping real packages to one Component
	ComponentName = string

	// abstract VendorName useful for mapping real vendor packages to one Vendor
	VendorName = string

	Document interface {
		Reference() models.Reference

		// Spec Version
		Version() speca.ReferableInt

		// Spec relative WorkingDirectory to root, prepend this to all path's from spec
		WorkingDirectory() speca.ReferableString

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
		IsDependOnAnyVendor() speca.ReferableBool
	}

	ExcludedDirectories interface {
		Reference() models.Reference

		// list of directories
		// examples:
		// 	- internal/test
		//	- vendor
		//	- .idea
		List() []speca.ReferableString
	}

	ExcludedFilesRegExp interface {
		Reference() models.Reference

		// list of regexp's
		// examples:
		// 	- "^.*_test\\.go$"
		List() []speca.ReferableString
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
		ImportPath() speca.ReferableString
	}

	CommonVendors interface {
		Reference() models.Reference

		// List of Vendors that can by imported to any project package
		List() []speca.ReferableString
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
		LocalPath() speca.ReferableString
	}

	CommonComponents interface {
		Reference() models.Reference

		// List of Components that can by imported to any project package
		List() []speca.ReferableString
	}

	Dependencies interface {
		Reference() models.Reference

		// Dependencies map between Components and DependencyRule`s
		Map() map[ComponentName]DependencyRule
	}

	DependencyRule interface {
		Reference() models.Reference

		// List of Component names, that can by imported to described component
		MayDependOn() []speca.ReferableString

		// List of Vendor names, that can by imported to described component
		CanUse() []speca.ReferableString

		// described component can import any other local namespace packages
		AnyProjectDeps() speca.ReferableBool

		// described component can import any other vendor namespace packages
		AnyVendorDeps() speca.ReferableBool
	}
)
