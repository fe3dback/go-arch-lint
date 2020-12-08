package project

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	Resolver struct {
	}

	resolveContext struct {
		projectDirectory    string
		moduleName          string
		excludePaths        []models.ResolvedPath
		excludeFileMatchers []*regexp.Regexp

		tokenSet *token.FileSet
		results  []models.ResolvedFile
	}
)

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Resolve(
	projectDirectory string,
	moduleName string,
	excludePaths []models.ResolvedPath,
	excludeFileMatchers []*regexp.Regexp,
) ([]models.ResolvedFile, error) {
	ctx := resolveContext{
		projectDirectory:    projectDirectory,
		moduleName:          moduleName,
		excludePaths:        excludePaths,
		excludeFileMatchers: excludeFileMatchers,

		tokenSet: token.NewFileSet(),
		results:  []models.ResolvedFile{},
	}

	err := filepath.Walk(ctx.projectDirectory, func(path string, info os.FileInfo, err error) error {
		return resolveFile(&ctx, path, info, err)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk project tree: %w", err)
	}

	return ctx.results, nil
}

func resolveFile(ctx *resolveContext, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || !inScope(ctx, path) {
		return nil
	}

	return parse(ctx, path)
}

func inScope(ctx *resolveContext, path string) bool {
	if filepath.Ext(path) != ".go" {
		return false
	}

	for _, excludePath := range ctx.excludePaths {
		if strings.HasPrefix(path, excludePath.AbsPath) {
			return false
		}
	}

	for _, matcher := range ctx.excludeFileMatchers {
		if matcher.Match([]byte(path)) {
			return false
		}
	}

	return true
}

func parse(ctx *resolveContext, path string) error {
	fileAst, err := parser.ParseFile(ctx.tokenSet, path, nil, parser.ImportsOnly)
	if err != nil {
		return fmt.Errorf("failed to parse go source code at '%s': %v", path, err)
	}

	ctx.results = append(ctx.results, models.ResolvedFile{
		Path:    path,
		Imports: extractImports(ctx, fileAst),
	})

	return nil
}

func extractImports(ctx *resolveContext, fileAst *ast.File) []models.ResolvedImport {
	imports := make([]models.ResolvedImport, 0)

	for _, goImport := range fileAst.Imports {
		importPath := strings.Trim(goImport.Path.Value, "\"")
		imports = append(imports, models.ResolvedImport{
			Name:       importPath,
			ImportType: getImportType(ctx, importPath),
		})
	}

	return imports
}

func getImportType(ctx *resolveContext, importPath string) models.ImportType {
	if !strings.Contains(importPath, ".") {
		return models.ImportTypeStdLib
	}

	if strings.HasPrefix(importPath, ctx.moduleName) {
		return models.ImportTypeProject
	}

	return models.ImportTypeVendor
}
