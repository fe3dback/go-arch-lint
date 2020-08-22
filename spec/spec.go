package spec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/goccy/go-yaml"
	structvalidator "gopkg.in/validator.v2"
)

const supportedVersion = 1

type (
	yamlVendorName       = string
	yamlComponentName    = string
	yamlExcludeLocalPath = string

	YamlSpec struct {
		Version            int                                 `yaml:"version"`
		Allow              YamlAllow                           `yaml:"allow"`
		Vendors            map[yamlVendorName]YamlVendor       `yaml:"vendors"`
		Exclude            []yamlExcludeLocalPath              `yaml:"exclude"`
		ExcludeFilesRegExp []string                            `yaml:"excludeFiles"`
		Components         map[yamlComponentName]YamlComponent `yaml:"components"`
		Dependencies       map[yamlComponentName]YamlRules     `yaml:"deps"`
		CommonComponents   []yamlComponentName                 `yaml:"commonComponents"`
		CommonVendors      []yamlVendorName                    `yaml:"commonVendors"`
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
		MayDependOn    []yamlComponentName `yaml:"mayDependOn"`
		CanUse         []yamlVendorName    `yaml:"canUse"`
		AnyProjectDeps bool                `yaml:"anyProjectDeps"`
		AnyVendorDeps  bool                `yaml:"anyVendorDeps"`
	}

	yamlSpecValidator struct {
		spec          YamlSpec
		rootDirectory string
	}
)

func newSpec(archFile string, rootDirectory string) (YamlSpec, error) {
	spec := YamlSpec{}

	data, err := ioutil.ReadFile(archFile)
	if err != nil {
		return spec, fmt.Errorf("can`t open '%s': %v", archFile, err)
	}

	reader := bytes.NewBuffer(data)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	err = decoder.Decode(&spec)
	if err != nil {
		return spec, fmt.Errorf("can`t parse yaml in '%s': %v", archFile, err)
	}

	specValidator := newValidator(spec, data, rootDirectory)
	err = specValidator.validate()
	if err != nil {
		return spec, fmt.Errorf("spec '%s' invalid: %v", archFile, err)
	}

	return spec, nil
}

func (sv *yamlSpecValidator) validateStruct(spec YamlSpec) error {
	var validators = map[string]func(interface{}, string) error{
		"xVendorPath":              sv.vendorPath,
		"xKnownComponent":          sv.validateIsKnownComponentName,
		"xKnownComponentSlice":     sv.validateIsKnownComponentNameSlice,
		"xMapNameIsKnownComponent": sv.validateMapNameIsKnownComponentName,
	}

	for name, fn := range validators {
		err := structvalidator.SetValidationFunc(name, fn)
		if err != nil {
			return fmt.Errorf("failed to create validator '%s': %v", name, err)
		}
	}

	err := structvalidator.Validate(spec)
	if err != nil {
		return fmt.Errorf("archFile invalid: %v", err)
	}

	return nil
}

func (sv *yamlSpecValidator) vendorPath(value interface{}, _ string) error {
	path, ok := value.(string)
	if !ok {
		return fmt.Errorf("should by string field")
	}

	return sv.checkDirectory(fmt.Sprintf("%s/vendor/%s", sv.rootDirectory, path))
}

func (sv *yamlSpecValidator) validateIsKnownComponentName(value interface{}, _ string) error {
	name, ok := value.(string)
	if !ok {
		return fmt.Errorf("should by string")
	}

	for componentName := range sv.spec.Components {
		if name == componentName {
			return nil
		}
	}

	return fmt.Errorf("component '%s' not known", name)
}

func (sv *yamlSpecValidator) validateIsKnownComponentNameSlice(value interface{}, param string) error {
	names, ok := value.([]string)
	if !ok {
		return fmt.Errorf("should by list of string's")
	}

	for _, name := range names {
		if err := sv.validateIsKnownComponentName(name, param); err != nil {
			return fmt.Errorf("invalid component '%s': %v", name, err)
		}
	}

	return nil
}

func (sv *yamlSpecValidator) validateMapNameIsKnownComponentName(value interface{}, param string) error {
	names := reflect.ValueOf(value).MapKeys()
	for _, name := range names {
		if err := sv.validateIsKnownComponentName(name.String(), param); err != nil {
			return fmt.Errorf("invalid component '%s': %v", name.String(), err)
		}
	}

	return nil
}

func (sv *yamlSpecValidator) checkDirectory(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' not exist", path)
		}
	}

	return nil
}
