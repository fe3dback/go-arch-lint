package assembler

import (
	"fmt"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Assembler struct {
	projectInfoFetcher projectInfoFetcher
	pathHelper         pathHelper
}

func NewAssembler(projectInfoFetcher projectInfoFetcher, pathHelper pathHelper) *Assembler {
	return &Assembler{
		projectInfoFetcher: projectInfoFetcher,
		pathHelper:         pathHelper,
	}
}

func (a *Assembler) Assemble(conf models.Config) (spec models.Spec, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed assemble: %v", r)
		}
	}()

	projectInfo, err := a.projectInfoFetcher.Fetch()
	if err != nil {
		return models.Spec{}, fmt.Errorf("failed fetch project info: %w", err)
	}

	components := make([]models.SpecComponent, 0, conf.Dependencies.Map.Len())
	conf.Dependencies.Map.Each(func(name models.ComponentName, rules models.ConfigComponentDependencies, reference models.Reference) {
		definition, definitionRef, exist := conf.Components.Map.Get(name)
		if !exist {
			panic(fmt.Errorf("component '%s' not found", name))
		}

		var deepScan models.Ref[bool]
		if rules.DeepScan.Defined {
			deepScan = rules.DeepScan.Value
		} else {
			deepScan = conf.Settings.DeepScan
		}

		tagsAllowedAll, tagsAllowedWhiteList := a.figureOutAllowedStructTags(&conf, &rules)

		matchedFiles, err := a.findOwnedFiles(conf.WorkingDirectory.Value, definition)
		if err != nil {
			panic(fmt.Errorf("failed finding owned files by component '%s': %w", name, err))
		}

		matchedPackages := a.extractUniquePackages(matchedFiles)

		component := models.SpecComponent{
			Name:                models.NewRef(name, definitionRef),
			DefinitionComponent: definitionRef,
			DefinitionDeps:      reference,
			DeepScan:            deepScan,
			StrictMode:          conf.Settings.Imports.StrictMode,
			AllowAllProjectDeps: rules.AnyProjectDeps,
			AllowAllVendorDeps:  rules.AnyVendorDeps,
			AllowAllTags:        tagsAllowedAll,
			AllowedTags:         tagsAllowedWhiteList,
			MayDependOn:         rules.MayDependOn,
			CanUse:              rules.CanUse,
			MatchedFiles:        matchedFiles,
			MatchedPackages:     matchedPackages,
		}

		components = append(components, component)
	})

	return models.Spec{
		Project:          projectInfo,
		WorkingDirectory: conf.WorkingDirectory,
		Components:       components,
	}, nil
}

func (a *Assembler) figureOutAllowedStructTags(conf *models.Config, rules *models.ConfigComponentDependencies) (models.Ref[bool], models.RefSlice[models.StructTag]) {
	if conf.Settings.Tags.Allowed.Value == models.ConfigSettingsTagsEnumAll {
		return models.NewRef(true, conf.Settings.Tags.Allowed.Ref), nil
	}

	globalTags := conf.Settings.Tags.AllowedList
	localTags := rules.CanContainTags

	allowedList := make(models.RefSlice[models.StructTag], 0, len(globalTags)+len(localTags))
	allowedList = append(allowedList, globalTags...)
	allowedList = append(allowedList, localTags...)

	return models.NewRef(false, conf.Settings.Tags.Allowed.Ref), allowedList
}

func (a *Assembler) findOwnedFiles(workingDirectory models.PathRelative, component models.ConfigComponent) ([]models.PathRelative, error) {
	filePaths := make([]models.PathRelative, 0, 32)

	for _, globPath := range component.In {
		files, err := a.pathHelper.FindProjectFiles(models.FileQuery{
			Path:               globPath.Value,
			WorkingDirectory:   workingDirectory,
			Type:               models.FileMatchQueryTypeOnlyFiles,
			ExcludeDirectories: nil, // todo
			ExcludeFiles:       nil, // todo
			ExcludeRegexp:      nil, // todo
			Extensions:         []string{"go"},
		})

		if err != nil {
			return nil, fmt.Errorf("matching glob path failed '%v': %w", globPath.Value, err)
		}

		for _, file := range files {
			filePaths = append(filePaths, file.PathRel)
		}
	}

	return filePaths, nil
}

func (a *Assembler) extractUniquePackages(files []models.PathRelative) []models.PathRelative {
	packages := make([]models.PathRelative, 0, len(files))
	unique := make(map[models.PathRelative]any)

	for _, file := range files {
		packagePath := models.PathRelative(filepath.Dir(string(file)))

		if _, ok := unique[packagePath]; ok {
			continue
		}

		unique[packagePath] = struct{}{}
		packages = append(packages, packagePath)
	}

	return packages
}
