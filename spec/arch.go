package spec

import (
	"fmt"
	"regexp"
	"strings"

	pathresolv "github.com/fe3dback/go-arch-lint/path"
)

type (
	Arch struct {
		rootDirectory string
		moduleName    string

		Allow               YamlAllow
		Components          []*Component
		Exclude             []*ResolvedPath
		ExcludeFilesMatcher []*regexp.Regexp
	}

	ComponentName = string
	Component     struct {
		Name           ComponentName
		LocalPathMask  string
		ResolvedPaths  []*ResolvedPath
		MayDependOn    []ComponentName
		AllowedImports []*ResolvedPath
		SpecialFlags   *SpecialFlags
	}

	ResolvedPath struct {
		ImportPath string
		LocalPath  string
		AbsPath    string
	}

	SpecialFlags struct {
		AllowAllProjectDeps bool
		AllowAllVendorDeps  bool
	}
)

func NewArch(archFile string, moduleName string, rootDirectory string) (*Arch, error) {
	spec, err := newSpec(archFile, rootDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to parse archfile: %v", err)
	}

	arch := Arch{
		rootDirectory:       rootDirectory,
		moduleName:          moduleName,
		Allow:               spec.Allow,
		Components:          make([]*Component, 0),
		Exclude:             make([]*ResolvedPath, 0),
		ExcludeFilesMatcher: make([]*regexp.Regexp, 0),
	}

	err = arch.assembleComponents(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble arch components: %v", err)
	}

	err = arch.assembleExclude(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble arch components: %v", err)
	}

	err = arch.assembleExcludeFilesMatcher(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble arch excludeFiles: %v", err)
	}

	return &arch, nil
}

func (a *Arch) assembleComponents(spec YamlSpec) error {
	for yamlName, yamlComponent := range spec.Components {
		depMeta := spec.Dependencies[yamlName]

		mayDependOn := make([]ComponentName, 0)
		for _, name := range depMeta.MayDependOn {
			mayDependOn = append(mayDependOn, name)
		}

		resolvedPath, err := a.assembleResolvedPaths(yamlComponent.LocalPath)
		if err != nil {
			return fmt.Errorf("failed to assemble component path's: %v", err)
		}

		allowedImports, err := a.assembleAllowedImports(spec, mayDependOn)
		if err != nil {
			return fmt.Errorf("failed to assemble component path's: %v", err)
		}

		a.Components = append(a.Components, &Component{
			Name:           yamlName,
			LocalPathMask:  yamlComponent.LocalPath,
			ResolvedPaths:  resolvedPath,
			MayDependOn:    mayDependOn,
			AllowedImports: allowedImports,
			SpecialFlags: &SpecialFlags{
				AllowAllProjectDeps: depMeta.AnyProjectDeps,
				AllowAllVendorDeps:  depMeta.anyVendorDeps,
			},
		})
	}

	return nil
}

func (a *Arch) assembleExclude(spec YamlSpec) error {
	for _, yamlRelativePath := range spec.Exclude {
		resolvedPath, err := a.assembleResolvedPaths(yamlRelativePath)
		if err != nil {
			return fmt.Errorf("failed to assemble exclude '%s' path's: %v", yamlRelativePath, err)
		}

		a.Exclude = append(a.Exclude, resolvedPath...)
	}

	return nil
}

func (a *Arch) assembleExcludeFilesMatcher(spec YamlSpec) error {
	for _, regString := range spec.ExcludeFilesRegExp {
		matcher, err := regexp.Compile(regString)
		if err != nil {
			return fmt.Errorf("failed to compile regular expression '%s': %v", regString, err)
		}

		a.ExcludeFilesMatcher = append(a.ExcludeFilesMatcher, matcher)
	}

	return nil
}

func (a *Arch) assembleResolvedPaths(localPathMask string) ([]*ResolvedPath, error) {
	list := make([]*ResolvedPath, 0)

	absPath := fmt.Sprintf("%s/%s", a.rootDirectory, localPathMask)
	resolved, err := pathresolv.ResolvePath(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path '%s'", absPath)
	}

	for _, absResolvedPath := range resolved {
		localPath := strings.TrimPrefix(absResolvedPath, fmt.Sprintf("%s/", a.rootDirectory))
		localPath = strings.TrimRight(localPath, "/")
		importPath := fmt.Sprintf("%s/%s", a.moduleName, localPath)

		list = append(list, &ResolvedPath{
			ImportPath: importPath,
			LocalPath:  localPath,
			AbsPath:    absResolvedPath,
		})
	}

	return list, nil
}

func (a *Arch) assembleAllowedImports(spec YamlSpec, componentNames []ComponentName) ([]*ResolvedPath, error) {
	list := make([]*ResolvedPath, 0)

	for _, name := range componentNames {
		maskPath := spec.Components[name].LocalPath

		resolved, err := a.assembleResolvedPaths(maskPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve mask '%s'", maskPath)
		}

		for _, resolvedPath := range resolved {
			list = append(list, resolvedPath)
		}
	}

	return list, nil
}
