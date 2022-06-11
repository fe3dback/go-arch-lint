package checker

import (
	"fmt"
	"path"
	"regexp"
	"sync"

	"github.com/fe3dback/go-arch-lint/internal/glue/deepscan"
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type DeepScan struct {
	scanner *deepscan.Searcher
	spec    speca.Spec
	result  models.CheckResult

	sync.Mutex
}

func NewDeepScan() *DeepScan {
	return &DeepScan{
		scanner: deepscan.NewSearcher(),
	}
}

func (c *DeepScan) Check(spec speca.Spec) (models.CheckResult, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.spec = spec
	c.result = models.CheckResult{}

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

		err := c.checkGate(cmp, usage, &gate)
		if err != nil {
			return fmt.Errorf("failed check gate '%s': %w",
				gate.Definition.Place,
				err,
			)
		}
	}

	return nil
}

func (c *DeepScan) checkGate(cmp *speca.Component, usage *deepscan.InjectionMethod, gate *deepscan.Gate) error {
	for _, implementation := range gate.Implementations {
		err := c.checkImplementation(cmp, usage, gate, &implementation)
		if err != nil {
			return fmt.Errorf("failed check implementation '%s': %w",
				implementation.Injector.Definition,
				err,
			)
		}
	}

	return nil
}

func (c *DeepScan) checkImplementation(
	cmp *speca.Component,
	usage *deepscan.InjectionMethod,
	gate *deepscan.Gate,
	imp *deepscan.Implementation,
) error {
	injectedImport := imp.Target.Definition.Import

	for _, allowedImport := range cmp.AllowedProjectImports {
		if allowedImport.Value().ImportPath == injectedImport {
			return nil
		}
	}

	warn := models.CheckDeepscanWarning{
		Gate: models.CheckDeepscanWarningGate{
			ComponentName:     cmp.Name.Value(),
			MethodName:        gate.Name, // todo: public func name
			Definition:        c.definitionToReference(gate.Definition.Place),
			SourceCodePreview: nil, // todo, custom height = (offsetOfParam - offsetOfDefinition)+1
		},
		Dependency: models.CheckDeepscanWarningDependency{
			ComponentName: "todo", // todo
			Name: fmt.Sprintf("%s.%s",
				imp.Target.Definition.Pkg,
				imp.Target.StructName,
			),
			InjectionAST:      imp.Injector.CodeName,
			Injection:         c.definitionToReference(imp.Injector.Definition.Place),
			SourceCodePreview: nil, // todo, custom height = (offsetOfParam - offsetOfDefinition)+1
		},
		LineArt: models.CheckDeepscanWarningArt{ // todo
			ToRight:    false,
			OutPos:     0,
			InPos:      0,
			LineLength: 0,
		},
	}

	c.result.DeepscanWarnings = append(c.result.DeepscanWarnings, warn)
	return nil
}

func (c *DeepScan) definitionToReference(source deepscan.Position) models.Reference {
	return models.Reference{
		Valid:  true,
		File:   source.Filename,
		Line:   source.Line,
		Offset: source.Offset,
	}
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
