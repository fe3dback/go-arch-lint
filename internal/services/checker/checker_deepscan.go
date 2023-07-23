package checker

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/services/checker/deepscan"
	"golang.org/x/sync/errgroup"
)

type DeepScan struct {
	projectFilesResolver projectFilesResolver
	sourceCodeRenderer   sourceCodeRenderer

	scanner           *deepscan.Searcher
	spec              arch.Spec
	result            models.CheckResult
	fileComponents    map[string]string
	packageComponents map[string]string

	sync.Mutex
}

func NewDeepScan(projectFilesResolver projectFilesResolver, sourceCodeRenderer sourceCodeRenderer) *DeepScan {
	return &DeepScan{
		projectFilesResolver: projectFilesResolver,
		sourceCodeRenderer:   sourceCodeRenderer,
		scanner:              deepscan.NewSearcher(),
	}
}

// How cpus available vs workersCount
// try to utilize CPU, but not all, user
// can do other work, and linter is background process
// 1 = 1   4 = 3   7 = 5
// 2 = 2   5 = 4   8 = 6
// 3 = 2   6 = 4   ...
func (c *DeepScan) workersCount() int {
	// currently scan algorithm can`t do work in ||
	// because of mux locks
	// its working, but 8 workers will scan with same speed that have 1
	//
	// todo: adapt scan algorithm to concurrent processing
	// todo: optimize scan speed
	return 1

	// max := runtime.NumCPU()
	// if max == 1 {
	// 	return 1
	// }
	// if max == 2 {
	// 	return 2
	// }
	//
	// half := int(math.Floor(float64(max) / 1.25))
	// if half < 2 {
	// 	half = 2
	// }
	//
	// return half
}

func (c *DeepScan) Check(ctx context.Context, spec arch.Spec) (models.CheckResult, error) {
	maxWorkers := c.workersCount()

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// -- prepare shared objects
	c.spec = spec
	c.result = models.CheckResult{}

	// -- prepare mapping file -> component
	mapping, err := c.projectFilesResolver.ProjectFiles(ctx, spec)
	if err != nil {
		return models.CheckResult{}, fmt.Errorf("failed resolve project files: %w", err)
	}

	c.fileComponents = map[string]string{}
	c.packageComponents = map[string]string{}

	for _, hold := range mapping {
		if hold.ComponentID == nil {
			continue
		}

		// cache file -> component ref
		c.fileComponents[hold.File.Path] = *hold.ComponentID

		// cache package -> component ref
		packagePath := filepath.Dir(hold.File.Path)
		c.packageComponents[packagePath] = *hold.ComponentID
	}

	// -- scan project
	pool := make(chan struct{}, maxWorkers)
	var wg errgroup.Group

	for _, component := range spec.Components {
		component := component

		pool <- struct{}{}
		wg.Go(func() error {
			defer func() {
				<-pool
			}()

			if component.DeepScan.Value != true {
				return nil
			}

			err := c.checkComponent(ctx, component)
			if err != nil {
				return fmt.Errorf("component '%s' check failed: %w",
					component.Name.Value,
					err,
				)
			}

			return nil
		})
	}

	err = wg.Wait()
	if err != nil {
		return models.CheckResult{}, err
	}

	return c.result, nil
}

func (c *DeepScan) checkComponent(ctx context.Context, cmp arch.Component) error {
	for _, packagePath := range cmp.ResolvedPaths {
		absPath := packagePath.Value.AbsPath
		matchedCmp, ok := c.packageComponents[absPath]
		if !ok {
			// component in excludes list
			continue
		}

		if matchedCmp != cmp.Name.Value {
			// this can be in cased of wildcard match, example:
			// cmp1: in: internal/code/common
			// cmp2: in: internal/code/**
			// and when cmp.Name == cmp2, this cmp2 still have "code/common" in resolvedPath's from cmp1
			// here we can skip this, because deps of cmp1 will be checked later.
			continue
		}

		err := c.scanPackage(ctx, &cmp, absPath)
		if err != nil {
			return fmt.Errorf("failed scan '%s': %w", absPath, err)
		}
	}

	return nil
}

func (c *DeepScan) scanPackage(ctx context.Context, cmp *arch.Component, absPackagePath string) error {
	usages, err := c.findUsages(ctx, absPackagePath)
	if err != nil {
		return fmt.Errorf("find usages failed: %w", err)
	}

	if len(usages) == 0 {
		return nil
	}

	for _, usage := range usages {
		err := c.checkUsage(ctx, cmp, &usage)
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

func (c *DeepScan) checkUsage(ctx context.Context, cmp *arch.Component, usage *deepscan.InjectionMethod) error {
	for _, gate := range usage.Gates {
		if len(gate.Implementations) == 0 {
			continue
		}

		err := c.checkGate(ctx, cmp, &gate)
		if err != nil {
			return fmt.Errorf("failed check gate '%s': %w",
				gate.ArgumentDefinition.Place,
				err,
			)
		}
	}

	return nil
}

func (c *DeepScan) checkGate(_ context.Context, cmp *arch.Component, gate *deepscan.Gate) error {
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
	cmp *arch.Component,
	gate *deepscan.Gate,
	imp *deepscan.Implementation,
) error {
	injectedImport := imp.Target.Definition.Import

	for _, allowedImport := range cmp.AllowedProjectImports {
		if allowedImport.Value.ImportPath == injectedImport {
			return nil
		}
	}

	targetPath := imp.Target.Definition.Place.File
	targetComponentID, targetDefined := c.fileComponents[targetPath]

	gatePath := gate.MethodDefinition.Place.File
	gateComponentID, gateDefined := c.fileComponents[gatePath]

	if !targetDefined || !gateDefined {
		// target component is vendor or std file, not described in mapping
		// we can check vendor libs too, but this requires another scan process

		// example of skipping target:
		// - $GOROOT/src/context/context.go (stdlib)
		// - /home/neo/go/src/example.com/ns/awesome/vendor/libs.example.com/good/producer/client.go (vendor)
		return nil
	}

	warn := models.CheckArchWarningDeepscan{
		Gate: models.DeepscanWarningGate{
			ComponentName: gateComponentID,
			MethodName:    gate.MethodName,
			RelativePath:  c.definitionToRelPath(gate.ArgumentDefinition.Place),
			Definition:    gate.ArgumentDefinition.Place,
		},
		Dependency: models.DeepscanWarningDependency{
			ComponentName: targetComponentID,
			Name: fmt.Sprintf("%s.%s",
				imp.Target.Definition.Pkg,
				imp.Target.StructName,
			),
			InjectionAST:  imp.Injector.CodeName,
			Injection:     imp.Injector.ParamDefinition.Place,
			InjectionPath: c.definitionToRelPath(imp.Injector.ParamDefinition.Place),
			SourceCodePreview: c.renderCode(
				imp.Injector.ParamDefinition.Place,
				imp.Injector.MethodDefinition.Place,
				imp.Injector.ParamDefinition.Place,
			),
		},
		Target: models.DeepscanWarningTarget{
			Definition:   imp.Target.Definition.Place,
			RelativePath: c.definitionToRelPath(imp.Target.Definition.Place),
		},
	}

	c.result.DeepscanWarnings = append(c.result.DeepscanWarnings, warn)
	return nil
}

func (c *DeepScan) renderCode(pointer, from, to common.Reference) []byte {
	return c.sourceCodeRenderer.SourceCode(
		common.NewReferenceRange(pointer.File, from.Line, pointer.Line, to.Line),
		false,
		false,
	)
}

func (c *DeepScan) sortPositions(a, b common.Reference) (min, max common.Reference) {
	if a.Line < b.Line {
		return a, b
	}

	return b, a
}

func (c *DeepScan) definitionToRelPath(source common.Reference) string {
	relativePath := strings.TrimPrefix(source.File, c.spec.RootDirectory.Value)
	return fmt.Sprintf("%s:%d", relativePath, source.Line)
}

func (c *DeepScan) findUsages(_ context.Context, absPackagePath string) ([]deepscan.InjectionMethod, error) {
	scanDirectory := path.Clean(fmt.Sprintf("%s/%s",
		c.spec.RootDirectory.Value,
		c.spec.WorkingDirectory.Value,
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

func (c *DeepScan) refPathToList(list []common.Referable[models.ResolvedPath]) []string {
	result := make([]string, 0)

	for _, refPath := range list {
		result = append(result, refPath.Value.AbsPath)
	}

	return result
}

func (c *DeepScan) refRegexpToList(list []common.Referable[*regexp.Regexp]) []*regexp.Regexp {
	result := make([]*regexp.Regexp, 0)

	for _, refPath := range list {
		result = append(result, refPath.Value)
	}

	return result
}
