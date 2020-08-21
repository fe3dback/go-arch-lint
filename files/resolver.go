package files

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const maxConcurrency = 16

const (
	ImportTypeStdLib ImportType = iota
	ImportTypeProject
	ImportTypeVendor
)

type (
	filePath = string
	chQueue  = chan filePath
	chError  = chan error
	chDone   = chan struct{}

	Resolver struct {
		projectDirectory    string
		moduleName          string
		excludePaths        []string
		excludeFileMatchers []*regexp.Regexp
		result              *ResolveResult
		tokenSet            *token.FileSet
		mux                 sync.Mutex
	}

	ImportType     uint8
	ResolvedImport struct {
		Name       string
		ImportType ImportType
	}

	ResolvedFile struct {
		Path    string
		Imports []ResolvedImport
	}

	ResolveResult struct {
		Files []*ResolvedFile
	}
)

func NewResolver(
	projectDirectory string,
	moduleName string,
	excludePaths []string,
	excludeFileMatchers []*regexp.Regexp,
) *Resolver {
	return &Resolver{
		projectDirectory:    projectDirectory,
		moduleName:          moduleName,
		excludePaths:        excludePaths,
		excludeFileMatchers: excludeFileMatchers,
		result:              &ResolveResult{},
		tokenSet:            token.NewFileSet(),
		mux:                 sync.Mutex{},
	}
}

func (r *Resolver) Resolve() (ResolveResult, error) {
	chDone := make(chDone)
	chErrors := make(chError)

	chQueue := r.walk(chErrors)
	r.handleQueue(chQueue, chErrors, chDone)

	select {
	case err := <-chErrors:
		return ResolveResult{}, err
	case <-chDone:
		return *r.result, nil
	}
}

func (r *Resolver) walk(chErr chError) chQueue {
	queue := make(chQueue, 40000)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		for _, excludePath := range r.excludePaths {
			if strings.HasPrefix(path, excludePath) {
				return nil
			}
		}

		if !r.isGoFile(path) {
			return nil
		}

		for _, matcher := range r.excludeFileMatchers {
			if matcher.Match([]byte(path)) {
				return nil
			}
		}

		queue <- path
		return nil
	}

	go func() {
		err := filepath.Walk(r.projectDirectory, walkFn)
		if err != nil {
			chErr <- fmt.Errorf("failed to walk project tree: %v", err)
		}

		close(queue)
	}()

	return queue
}

func (r *Resolver) isGoFile(path string) bool {
	return filepath.Ext(path) == ".go"
}

func (r *Resolver) handleQueue(queue chQueue, chErr chError, done chDone) {
	concurrency := make(chan struct{}, maxConcurrency)

	go func() {
		var wg sync.WaitGroup
		for path := range queue {
			concurrency <- struct{}{}
			go func(path filePath) {
				wg.Add(1)
				defer wg.Done()

				err := r.parse(path)
				if err != nil {
					chErr <- fmt.Errorf("failed to parse '%s': %v", path, err)
					return
				}

				<-concurrency
			}(path)
		}

		wg.Wait()
		done <- struct{}{}
	}()
}

func (r *Resolver) parse(path filePath) error {
	fileAst, err := parser.ParseFile(r.tokenSet, path, nil, parser.ImportsOnly)
	if err != nil {
		return fmt.Errorf("failed to parse go at '%s': %v", path, err)
	}

	imports := r.extractImports(fileAst)

	r.mux.Lock()
	r.result.Files = append(r.result.Files, &ResolvedFile{
		Path:    path,
		Imports: imports,
	})
	r.mux.Unlock()

	return nil
}

func (r *Resolver) extractImports(fileAst *ast.File) []ResolvedImport {
	imports := make([]ResolvedImport, 0)

	for _, goImport := range fileAst.Imports {
		importPath := strings.Trim(goImport.Path.Value, "\"")
		imports = append(imports, ResolvedImport{
			Name:       importPath,
			ImportType: r.getImportType(importPath),
		})
	}

	return imports
}

func (r *Resolver) getImportType(importPath string) ImportType {
	if !strings.Contains(importPath, ".") {
		return ImportTypeStdLib
	}

	if strings.HasPrefix(importPath, r.moduleName) {
		return ImportTypeProject
	}

	return ImportTypeVendor
}
