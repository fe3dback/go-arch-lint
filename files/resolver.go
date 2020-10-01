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

	"github.com/fe3dback/go-arch-lint/models"
)

type (
	Resolver struct {
		projectDirectory    string
		moduleName          string
		excludePaths        []string
		excludeFileMatchers []*regexp.Regexp
		resolvedFiles       []*models.ResolvedFile
		tokenSet            *token.FileSet
		mux                 sync.Mutex
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
		resolvedFiles:       make([]*models.ResolvedFile, 0),
		tokenSet:            token.NewFileSet(),
		mux:                 sync.Mutex{},
	}
}

func (r *Resolver) Resolve() ([]*models.ResolvedFile, error) {
	err := filepath.Walk(r.projectDirectory, r.resolveFile)
	if err != nil {
		return nil, fmt.Errorf("failed to walk project tree: %v", err)
	}

	return r.resolvedFiles, nil
}

func (r *Resolver) resolveFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || !r.inScope(path) {
		return nil
	}

	return r.parse(path)
}

func (r *Resolver) inScope(path string) bool {
	if filepath.Ext(path) != ".go" {
		return false
	}

	for _, excludePath := range r.excludePaths {
		if strings.HasPrefix(path, excludePath) {
			return false
		}
	}

	for _, matcher := range r.excludeFileMatchers {
		if matcher.Match([]byte(path)) {
			return false
		}
	}

	return true
}

func (r *Resolver) parse(path string) error {
	fileAst, err := parser.ParseFile(r.tokenSet, path, nil, parser.ImportsOnly)
	if err != nil {
		return fmt.Errorf("failed to parse go source code at '%s': %v", path, err)
	}

	imports := r.extractImports(fileAst)

	r.mux.Lock()
	r.resolvedFiles = append(r.resolvedFiles, &models.ResolvedFile{
		Path:    path,
		Imports: imports,
	})
	r.mux.Unlock()

	return nil
}

func (r *Resolver) extractImports(fileAst *ast.File) []models.ResolvedImport {
	imports := make([]models.ResolvedImport, 0)

	for _, goImport := range fileAst.Imports {
		importPath := strings.Trim(goImport.Path.Value, "\"")
		imports = append(imports, models.ResolvedImport{
			Name:       importPath,
			ImportType: r.getImportType(importPath),
		})
	}

	return imports
}

func (r *Resolver) getImportType(importPath string) models.ImportType {
	if !strings.Contains(importPath, ".") {
		return models.ImportTypeStdLib
	}

	if strings.HasPrefix(importPath, r.moduleName) {
		return models.ImportTypeProject
	}

	return models.ImportTypeVendor
}
