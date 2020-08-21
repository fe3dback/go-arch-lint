package spec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	pathresolv "github.com/fe3dback/go-arch-lint/path"

	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

const supportedVersion = 1

type (
	yamlVendorName       = string
	yamlComponentName    = string
	yamlExcludeLocalPath = string

	YamlSpec struct {
		Version            int                                 `yaml:"version" validate:"nonnil, min=1"`
		Allow              YamlAllow                           `yaml:"allow" validate:"nonnil"`
		Vendors            map[yamlVendorName]YamlVendor       `yaml:"vendors"`
		Exclude            []yamlExcludeLocalPath              `yaml:"exclude" validate:"xPathSlice"`
		ExcludeFilesRegExp []string                            `yaml:"excludeFiles"`
		Components         map[yamlComponentName]YamlComponent `yaml:"components" validate:"nonnil, xMapNameIsKnownComponent"`
		Dependencies       map[yamlComponentName]YamlRules     `yaml:"deps" validate:"nonnil, xMapNameIsKnownComponent"`
	}

	YamlAllow struct {
		DepOnAnyVendor bool `yaml:"depOnAnyVendor" validate:"nonnil"`
	}

	YamlVendor struct {
		ImportPath string `yaml:"in" validate:"min=1, xVendorPath"`
	}

	YamlComponent struct {
		LocalPath string `yaml:"in" validate:"min=1, xPath"`
	}

	YamlRules struct {
		MayDependOn    []yamlComponentName `yaml:"mayDependOn" validate:"xKnownComponentSlice"`
		AnyProjectDeps bool                `yaml:"anyProjectDeps"`
		anyVendorDeps  bool                `yaml:"anyVendorDeps"`
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

	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		return spec, fmt.Errorf("can`t parse yaml in '%s': %v", archFile, err)
	}

	specValidator := yamlSpecValidator{
		spec:          spec,
		rootDirectory: rootDirectory,
	}

	err = specValidator.validateSpec(spec)
	if err != nil {
		return spec, fmt.Errorf("spec '%s' invalid: %v", archFile, err)
	}

	return spec, nil
}

func (sv *yamlSpecValidator) validateSpec(spec YamlSpec) error {
	var validators = map[string]func(interface{}, string) error{
		"xVendorPath":              sv.vendorPath,
		"xPath":                    sv.validatePath,
		"xPathSlice":               sv.validatePathSlice,
		"xKnownComponent":          sv.validateIsKnownComponentName,
		"xKnownComponentSlice":     sv.validateIsKnownComponentNameSlice,
		"xMapNameIsKnownComponent": sv.validateMapNameIsKnownComponentName,
	}

	for name, fn := range validators {
		err := validator.SetValidationFunc(name, fn)
		if err != nil {
			return fmt.Errorf("failed to create validator '%s': %v", name, err)
		}
	}

	if spec.Version > supportedVersion {
		return fmt.Errorf("archFile has newer version %d, current support version: %d",
			spec.Version,
			supportedVersion,
		)
	}

	err := validator.Validate(spec)
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

func (sv *yamlSpecValidator) validatePath(value interface{}, _ string) error {
	localPath, ok := value.(string)
	if !ok {
		return fmt.Errorf("should by string field")
	}

	absPath := filepath.Clean(fmt.Sprintf("%s/%s", sv.rootDirectory, localPath))
	resolved, err := pathresolv.ResolvePath(absPath)
	if err != nil {
		return fmt.Errorf("failed to resolv path: %v", err)
	}

	return sv.checkDirectory(resolved...)
}

func (sv *yamlSpecValidator) validatePathSlice(value interface{}, param string) error {
	pathSlice, ok := value.([]string)
	if !ok {
		return fmt.Errorf("should by list of string's")
	}

	for _, path := range pathSlice {
		if err := sv.validatePath(path, param); err != nil {
			return fmt.Errorf("invalid path '%s': %v", path, err)
		}
	}

	return nil
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
