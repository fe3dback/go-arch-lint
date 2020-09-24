package archfile

const SupportedVersion = 1

type (
	YamlVendorName       = string
	YamlComponentName    = string
	YamlExcludeLocalPath = string

	YamlSpec struct {
		Version            int                                 `yaml:"version"`
		Allow              YamlAllow                           `yaml:"allow"`
		Vendors            map[YamlVendorName]YamlVendor       `yaml:"vendors"`
		Exclude            []YamlExcludeLocalPath              `yaml:"exclude"`
		ExcludeFilesRegExp []string                            `yaml:"excludeFiles"`
		Components         map[YamlComponentName]YamlComponent `yaml:"components"`
		Dependencies       map[YamlComponentName]YamlRules     `yaml:"deps"`
		CommonComponents   []YamlComponentName                 `yaml:"commonComponents"`
		CommonVendors      []YamlVendorName                    `yaml:"commonVendors"`
	}

	YamlAllow struct {
		DepOnAnyVendor bool `yaml:"depOnAnyVendor"`
	}

	YamlVendor struct {
		ImportPath string `yaml:"in"`
	}

	YamlComponent struct {
		LocalPath string `yaml:"in"`
	}

	YamlRules struct {
		MayDependOn    []YamlComponentName `yaml:"mayDependOn"`
		CanUse         []YamlVendorName    `yaml:"canUse"`
		AnyProjectDeps bool                `yaml:"anyProjectDeps"`
		AnyVendorDeps  bool                `yaml:"anyVendorDeps"`
	}
)
