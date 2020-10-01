package archfile

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
)

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

func NewYamlSpec(sourceCode []byte) (*YamlSpec, error) {
	reader := bytes.NewBuffer(sourceCode)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	spec := YamlSpec{}
	err := decoder.Decode(&spec)
	if err != nil {
		return nil, fmt.Errorf("can`t parse yaml: %w", err)
	}

	return &spec, nil
}
