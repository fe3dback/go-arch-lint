package spec

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	pathresolv "github.com/fe3dback/go-arch-lint/path"
	"github.com/fe3dback/go-arch-lint/spec/archfile"
)

type (
	Arch struct {
		rootDirectory string
		moduleName    string

		Allow               archfile.YamlAllow
		Components          []*Component
		Exclude             []*ResolvedPath
		ExcludeFilesMatcher []*regexp.Regexp
	}

	ComponentName = string
	VendorName    = string
	Component     struct {
		Name           ComponentName
		LocalPathMask  string
		ResolvedPaths  []*ResolvedPath
		MayDependOn    []ComponentName
		CanUse         []VendorName
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

func NewArch(archFile string, moduleName string, rootDirectory string) (*Arch, error, []YamlAnnotatedWarning) {
	spec, specParseErr := newSpec(archFile, rootDirectory)
	if specParseErr.Err != nil {
		return nil, fmt.Errorf("failed to parse archfile: %v", specParseErr.Err), specParseErr.Warnings
	}

	arch := Arch{
		rootDirectory:       rootDirectory,
		moduleName:          moduleName,
		Allow:               spec.Allow,
		Components:          make([]*Component, 0),
		Exclude:             make([]*ResolvedPath, 0),
		ExcludeFilesMatcher: make([]*regexp.Regexp, 0),
	}

	err := arch.assembleComponents(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble arch components: %v", err), nil
	}

	err = arch.assembleExclude(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble arch components: %v", err), nil
	}

	err = arch.assembleExcludeFilesMatcher(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble arch excludeFiles: %v", err), nil
	}

	return &arch, nil, nil
}

func (a *Arch) assembleComponents(spec archfile.YamlSpec) error {
	for yamlName, yamlComponent := range spec.Components {
		depMeta := spec.Dependencies[yamlName]

		mayDependOn := make([]ComponentName, 0)
		for _, name := range depMeta.MayDependOn {
			mayDependOn = append(mayDependOn, name)
		}

		canUse := make([]VendorName, 0)
		for _, name := range depMeta.CanUse {
			canUse = append(canUse, name)
		}

		resolvedPath, err := a.assembleResolvedPaths(yamlComponent.LocalPath)
		if err != nil {
			return fmt.Errorf("failed to assemble component path's: %v", err)
		}

		allowedImports, err := a.assembleAllowedImports(spec, mayDependOn, canUse)
		if err != nil {
			return fmt.Errorf("failed to assemble component path's: %v", err)
		}

		a.Components = append(a.Components, &Component{
			Name:           yamlName,
			LocalPathMask:  yamlComponent.LocalPath,
			ResolvedPaths:  resolvedPath,
			MayDependOn:    mayDependOn,
			CanUse:         canUse,
			AllowedImports: allowedImports,
			SpecialFlags: &SpecialFlags{
				AllowAllProjectDeps: depMeta.AnyProjectDeps,
				AllowAllVendorDeps:  depMeta.AnyVendorDeps,
			},
		})
	}

	return nil
}

func (a *Arch) assembleExclude(spec archfile.YamlSpec) error {
	for _, yamlRelativePath := range spec.Exclude {
		resolvedPath, err := a.assembleResolvedPaths(yamlRelativePath)
		if err != nil {
			return fmt.Errorf("failed to assemble exclude '%s' path's: %v", yamlRelativePath, err)
		}

		a.Exclude = append(a.Exclude, resolvedPath...)
	}

	return nil
}

func (a *Arch) assembleExcludeFilesMatcher(spec archfile.YamlSpec) error {
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
			ImportPath: strings.TrimRight(importPath, "/"),
			LocalPath:  strings.TrimRight(localPath, "/"),
			AbsPath:    filepath.Clean(strings.TrimRight(absResolvedPath, "/")),
		})
	}

	return list, nil
}

func (a *Arch) assembleAllowedImports(
	spec archfile.YamlSpec,
	componentNames []ComponentName,
	vendorNames []VendorName,
) ([]*ResolvedPath, error) {
	list := make([]*ResolvedPath, 0)

	allowedComponents := make([]ComponentName, 0)
	allowedComponents = append(allowedComponents, componentNames...)
	allowedComponents = append(allowedComponents, spec.CommonComponents...)

	allowedVendors := make([]VendorName, 0)
	allowedVendors = append(allowedVendors, vendorNames...)
	allowedVendors = append(allowedVendors, spec.CommonVendors...)

	for _, name := range allowedComponents {
		maskPath := spec.Components[name].LocalPath

		resolved, err := a.assembleResolvedPaths(maskPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve mask '%s'", maskPath)
		}

		for _, resolvedPath := range resolved {
			list = append(list, resolvedPath)
		}
	}

	for _, name := range allowedVendors {
		importPath := spec.Vendors[name].ImportPath
		localPath := fmt.Sprintf("vendor/%s", importPath)
		absPath := fmt.Sprintf("%s/%s", a.rootDirectory, localPath)

		list = append(list, &ResolvedPath{
			ImportPath: importPath,
			LocalPath:  localPath,
			AbsPath:    absPath,
		})
	}

	return list, nil
}
