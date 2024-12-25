package scanner

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	astUtil "github.com/fe3dback/go-arch-lint/internal/services/common/ast"
	"golang.org/x/tools/go/packages"
)

type (
	Scanner struct {
		stdPackages map[string]struct{}
	}

	resolveContext struct {
		projectDirectory    string
		moduleName          string
		excludePaths        []models.ResolvedPath
		excludeFileMatchers []*regexp.Regexp

		tokenSet *token.FileSet
		results  []models.ProjectFile
	}
)

func NewScanner() *Scanner {
	scanner := &Scanner{
		stdPackages: make(map[string]struct{}, 255),
	}

	stdPackages, err := packages.Load(nil, "std")
	if err != nil {
		panic(fmt.Errorf("failed load std packages"))
	}

	for _, stdPackage := range stdPackages {
		scanner.stdPackages[stdPackage.ID] = struct{}{}
	}

	return scanner
}

func (r *Scanner) Scan(
	_ context.Context,
	projectDirectory string,
	moduleName string,
	excludePaths []models.ResolvedPath,
	excludeFileMatchers []*regexp.Regexp,
) ([]models.ProjectFile, error) {
	rctx := resolveContext{
		projectDirectory:    projectDirectory,
		moduleName:          moduleName,
		excludePaths:        excludePaths,
		excludeFileMatchers: excludeFileMatchers,

		tokenSet: token.NewFileSet(),
		results:  []models.ProjectFile{},
	}

	err := filepath.Walk(rctx.projectDirectory, func(path string, info os.FileInfo, err error) error {
		return r.resolveFile(&rctx, path, info, err)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk project tree: %w", err)
	}

	return rctx.results, nil
}

func (r *Scanner) resolveFile(ctx *resolveContext, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || !r.inScope(ctx, path) {
		return nil
	}

	return r.parse(ctx, path)
}

func (r *Scanner) inScope(ctx *resolveContext, path string) bool {
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

func (r *Scanner) parse(ctx *resolveContext, path string) error {
	fileAst, err := parser.ParseFile(ctx.tokenSet, path, nil, parser.ImportsOnly)
	if err != nil {
		return fmt.Errorf("failed to parse go source code at '%s': %w", path, err)
	}

	ctx.results = append(ctx.results, models.ProjectFile{
		Path:    path,
		Imports: r.extractImports(ctx, fileAst),
	})

	return nil
}

func (r *Scanner) extractImports(ctx *resolveContext, fileAst *ast.File) []models.ResolvedImport {
	imports := make([]models.ResolvedImport, 0)

	for _, goImport := range fileAst.Imports {
		importPath := strings.Trim(goImport.Path.Value, "\"")
		imports = append(imports, models.ResolvedImport{
			Name:       importPath,
			ImportType: r.getImportType(ctx, importPath),
			Reference:  astUtil.PositionFromToken(ctx.tokenSet.Position(goImport.Pos())),
		})
	}

	return imports
}

func (r *Scanner) getImportType(ctx *resolveContext, importPath string) models.ImportType {
	if _, ok := r.stdPackages[importPath]; ok {
		return models.ImportTypeStdLib
	}

	// We can't use a straight prefix match here because the module name could be a substring of the import path.
	// For example, if the module name is "example.com/foo/bar", we do not want to match "example.com/foo/bar-utils"
	if importPath == ctx.moduleName || strings.HasPrefix(importPath, ctx.moduleName+"/") {
		return models.ImportTypeProject
	}

	return models.ImportTypeVendor
}
