package checker

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/fe3dback/go-arch-lint/internal/glue/deepscan"
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type DeepScan struct {
	projectFilesResolver projectFilesResolver
	sourceCodeRenderer   sourceCodeRenderer

	scanner        *deepscan.Searcher
	spec           speca.Spec
	result         models.CheckResult
	fileComponents map[string]string

	sync.Mutex
}

func NewDeepScan(projectFilesResolver projectFilesResolver, sourceCodeRenderer sourceCodeRenderer) *DeepScan {
	return &DeepScan{
		projectFilesResolver: projectFilesResolver,
		sourceCodeRenderer:   sourceCodeRenderer,
		scanner:              deepscan.NewSearcher(),
	}
}

func (c *DeepScan) Check(spec speca.Spec) (models.CheckResult, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// -- prepare shared objects
	c.spec = spec
	c.result = models.CheckResult{}

	// -- prepare mapping file -> component
	mapping, err := c.projectFilesResolver.ProjectFiles(spec)
	if err != nil {
		return models.CheckResult{}, fmt.Errorf("failed resolve project files: %w", err)
	}

	c.fileComponents = map[string]string{}
	for _, hold := range mapping {
		if hold.ComponentID == nil {
			continue
		}

		c.fileComponents[hold.File.Path] = *hold.ComponentID
	}

	// -- scan project
	for _, component := range spec.Components {
		if component.DeepScan.Value() != true {
			continue
		}

		err := c.checkComponent(component)
		if err != nil {
			return models.CheckResult{}, fmt.Errorf("component '%s' check failed: %w",
				component.Name.Value(),
				err,
			)
		}
	}

	return c.result, nil
}

func (c *DeepScan) checkComponent(cmp speca.Component) error {
	for _, packagePath := range cmp.ResolvedPaths {
		absPath := packagePath.Value().AbsPath
		err := c.scanPackage(&cmp, absPath)
		if err != nil {
			return fmt.Errorf("failed scan '%s': %w", absPath, err)
		}
	}

	return nil
}

func (c *DeepScan) scanPackage(cmp *speca.Component, absPackagePath string) error {
	usages, err := c.findUsages(absPackagePath)
	if err != nil {
		return fmt.Errorf("find usages failed: %w", err)
	}

	if len(usages) == 0 {
		return nil
	}

	for _, usage := range usages {
		err := c.checkUsage(cmp, &usage)
		if err != nil {
			return fmt.Errorf("failed check usage '%s' in '%s': %w",
				usage.Name,
				usage.Definition.Place,
				err,
			)
		}
	}

	return nil
}

func (c *DeepScan) checkUsage(cmp *speca.Component, usage *deepscan.InjectionMethod) error {
	for _, gate := range usage.Gates {
		if len(gate.Implementations) == 0 {
			continue
		}

		err := c.checkGate(cmp, &gate)
		if err != nil {
			return fmt.Errorf("failed check gate '%s': %w",
				gate.ArgumentDefinition.Place,
				err,
			)
		}
	}

	return nil
}

func (c *DeepScan) checkGate(cmp *speca.Component, gate *deepscan.Gate) error {
	for _, implementation := range gate.Implementations {
		err := c.checkImplementation(cmp, gate, &implementation)
		if err != nil {
			return fmt.Errorf("failed check implementation '%s': %w",
				implementation.Injector.ParamDefinition,
				err,
			)
		}
	}

	return nil
}

func (c *DeepScan) checkImplementation(
	cmp *speca.Component,
	gate *deepscan.Gate,
	imp *deepscan.Implementation,
) error {
	injectedImport := imp.Target.Definition.Import

	for _, allowedImport := range cmp.AllowedProjectImports {
		if allowedImport.Value().ImportPath == injectedImport {
			return nil
		}
	}

	targetName := imp.Target.Definition.Place.Filename
	targetComponentID, targetDefined := c.fileComponents[targetName]

	if !targetDefined {
		// target component not described in go-arch-lint config
		// so skip this warning, because linter show another warning
		// anyway that this target is not mapped
		return nil
	}

	warn := models.CheckArchWarningDeepscan{
		Gate: models.DeepscanWarningGate{
			ComponentName: cmp.Name.Value(),
			MethodName:    gate.MethodName,
			RelativePath:  c.definitionToRelPath(gate.ArgumentDefinition.Place),
			Definition:    c.definitionToReference(gate.ArgumentDefinition.Place),
		},
		Dependency: models.DeepscanWarningDependency{
			ComponentName: targetComponentID,
			Name: fmt.Sprintf("%s.%s",
				imp.Target.Definition.Pkg,
				imp.Target.StructName,
			),
			InjectionAST:  imp.Injector.CodeName,
			Injection:     c.definitionToReference(imp.Injector.ParamDefinition.Place),
			InjectionPath: c.definitionToRelPath(imp.Injector.ParamDefinition.Place),
			SourceCodePreview: c.renderCodeBetween(
				imp.Injector.MethodDefinition.Place,
				imp.Injector.ParamDefinition.Place,
			),
		},
	}

	c.result.DeepscanWarnings = append(c.result.DeepscanWarnings, warn)
	return nil
}

func (c *DeepScan) renderCodeBetween(from, to deepscan.Position) []byte {
	const maxCodeBlockHeight = 5

	if from.Filename != to.Filename {
		// invalid references
		return nil
	}

	if from.Line == to.Line {
		// its same line, render only it
		return c.renderCodeFrom(from, 1)
	}

	min, max := c.sortPositions(from, to)
	height := (max.Line - min.Line) + 1

	if height <= maxCodeBlockHeight {
		// render all block
		return c.renderCodeFrom(min, height)
	}

	// render only last useful line
	return c.renderCodeFrom(max, 1)
}

func (c *DeepScan) renderCodeFrom(from deepscan.Position, height int) []byte {
	const highlight = true
	return c.sourceCodeRenderer.SourceCodeWithoutOffset(
		c.definitionToReference(from),
		height,
		highlight,
	)
}

func (c *DeepScan) sortPositions(a, b deepscan.Position) (min, max deepscan.Position) {
	if a.Line < b.Line {
		return a, b
	}

	return b, a
}

func (c *DeepScan) definitionToReference(source deepscan.Position) models.Reference {
	return models.Reference{
		Valid:  true,
		File:   source.Filename,
		Line:   source.Line,
		Offset: source.Column,
	}
}

func (c *DeepScan) definitionToRelPath(source deepscan.Position) string {
	relativePath := strings.TrimPrefix(source.Filename, c.spec.RootDirectory.Value())
	return fmt.Sprintf("%s:%d", relativePath, source.Line)
}

func (c *DeepScan) findUsages(absPackagePath string) ([]deepscan.InjectionMethod, error) {
	scanDirectory := path.Clean(fmt.Sprintf("%s/%s",
		c.spec.RootDirectory.Value(),
		c.spec.WorkingDirectory.Value(),
	))
	excludeDirectories := c.refPathToList(c.spec.Exclude)
	excludeMatchers := c.refRegexpToList(c.spec.ExcludeFilesMatcher)

	criteria, err := deepscan.NewCriteria(
		deepscan.WithPackagePath(absPackagePath),
		deepscan.WithAnalyseScope(scanDirectory),
		deepscan.WithExcludedPath(excludeDirectories),
		deepscan.WithExcludedFileMatchers(excludeMatchers),
	)
	if err != nil {
		return nil, fmt.Errorf("failed prepare scan criteria: %w", err)
	}

	usages, err := c.scanner.Usages(criteria)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return usages, nil
}

func (c *DeepScan) refPathToList(list []speca.Referable[models.ResolvedPath]) []string {
	result := make([]string, 0)

	for _, refPath := range list {
		result = append(result, refPath.Value().AbsPath)
	}

	return result
}

func (c *DeepScan) refRegexpToList(list []speca.Referable[*regexp.Regexp]) []*regexp.Regexp {
	result := make([]*regexp.Regexp, 0)

	for _, refPath := range list {
		result = append(result, refPath.Value())
	}

	return result
}
