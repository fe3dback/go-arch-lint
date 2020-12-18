package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	VendorName       = string
	ComponentName    = string
	ExcludeLocalPath = string

	ArchDocument struct {
		Document  Document
		Integrity []speca.Notice
	}

	Document struct {
		Version            int                         `yaml:"version"`
		Allow              Allow                       `yaml:"allow"`
		Vendors            map[VendorName]Vendor       `yaml:"vendors"`
		Exclude            []ExcludeLocalPath          `yaml:"exclude"`
		ExcludeFilesRegExp []string                    `yaml:"excludeFiles"`
		Components         map[ComponentName]Component `yaml:"components"`
		Dependencies       map[ComponentName]Rules     `yaml:"deps"`
		CommonComponents   []ComponentName             `yaml:"commonComponents"`
		CommonVendors      []VendorName                `yaml:"commonVendors"`
	}

	Allow struct {
		DepOnAnyVendor bool `yaml:"depOnAnyVendor"`
	}

	Vendor struct {
		ImportPath string `yaml:"in"`
	}

	Component struct {
		LocalPath string `yaml:"in"`
	}

	Rules struct {
		MayDependOn    []ComponentName `yaml:"mayDependOn"`
		CanUse         []VendorName    `yaml:"canUse"`
		AnyProjectDeps bool            `yaml:"anyProjectDeps"`
		AnyVendorDeps  bool            `yaml:"anyVendorDeps"`
	}
)
