package assembler

import (
	"fmt"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/pkg/pathsort"
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

// nolint
func (a *Assembler) Assemble(conf models.Config) (spec models.Spec, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed assemble: %v\n%s", r, debug.Stack())
		}
	}()

	projectInfo, err := a.projectInfoFetcher.Fetch()
	if err != nil {
		return models.Spec{}, fmt.Errorf("failed fetch project info: %w", err)
	}

	components := make([]*models.SpecComponent, 0, conf.Dependencies.Map.Len())

	conf.Components.Map.Each(func(name models.ComponentName, definition models.ConfigComponent, definitionRef models.Reference) {

		rules, rulesRef, exist := conf.Dependencies.Map.Get(name)
		if !exist {
			// defaults for rules
			rules = models.ConfigComponentDependencies{}
			rulesRef = models.NewInvalidReference()
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

		components = append(components, &models.SpecComponent{
			Name:                models.NewRef(name, definitionRef),
			DefinitionComponent: definitionRef,
			DefinitionDeps:      rulesRef,
			DeepScan:            deepScan,
			StrictMode:          conf.Settings.Imports.StrictMode,
			AllowAllProjectDeps: rules.AnyProjectDeps,
			AllowAllVendorDeps:  rules.AnyVendorDeps,
			AllowAllTags:        tagsAllowedAll,
			AllowedTags:         tagsAllowedWhiteList,
			MayDependOn:         rules.MayDependOn,
			CanUse:              rules.CanUse,
			MatchPatterns:       definition.In,
			MatchedFiles:        matchedFiles,
		})
	})

	// copy matched files to owned files (but each file owned only by one component)
	a.calculateFilesOwnage(components)

	// find matched/owned packages from files
	for _, component := range components {
		component.MatchedPackages = a.extractUniquePackages(component.MatchedFiles)
		component.OwnedPackages = a.extractUniquePackages(component.OwnedFiles)
	}

	// sort paths
	for _, component := range components {
		pathsort.SortDescriptors(component.MatchedFiles)
		pathsort.SortDescriptors(component.MatchedPackages)
		pathsort.SortDescriptors(component.OwnedFiles)
		pathsort.SortDescriptors(component.OwnedPackages)
	}

	// finalize
	resultComponents := make([]models.SpecComponent, 0, len(components))
	for _, component := range components {
		resultComponents = append(resultComponents, *component)
	}

	return models.Spec{
		Project:          projectInfo,
		WorkingDirectory: conf.WorkingDirectory,
		Components:       resultComponents,
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

func (a *Assembler) findOwnedFiles(workingDirectory models.PathRelative, component models.ConfigComponent) ([]models.FileDescriptor, error) {
	list := make([]models.FileDescriptor, 0, 32)

	for _, globPath := range component.In {
		// convert directory glob to file scope glob
		fileGlob := globPath.Value

		if !strings.HasSuffix(string(fileGlob), "/**") {
			// rules:
			// "app"            | match only app itself (directory), but no files inside
			// "app/*"          | match all files inside app itself, but no directory and subdirs (and subdirs files)
			// "app/**"         | match all files inside app and all subdirs (will recursive files on any level)

			// convert "app" -> "app/*", because we want to find files inside this directory
			fileGlob = models.PathRelativeGlob(fmt.Sprintf("%s/*", fileGlob))
		}

		files, err := a.pathHelper.FindProjectFiles(models.FileQuery{
			Path:               fileGlob,
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

		list = append(list, files...)
	}

	return list, nil
}

func (a *Assembler) extractUniquePackages(files []models.FileDescriptor) []models.FileDescriptor {
	packages := make([]models.FileDescriptor, 0, len(files))
	unique := make(map[models.PathRelative]any)

	for _, file := range files {
		packagePathRelative := models.PathRelative(filepath.Dir(string(file.PathRel)))
		packagePathAbsolute := models.PathAbsolute(filepath.Dir(string(file.PathAbs)))

		if _, ok := unique[packagePathRelative]; ok {
			continue
		}

		unique[packagePathRelative] = struct{}{}
		packages = append(packages, models.FileDescriptor{
			PathRel: packagePathRelative,
			PathAbs: packagePathAbsolute,
			IsDir:   true,
		})
	}

	return packages
}
