package scanner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	Scanner struct {
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

var stdPackagesNames map[string]struct{}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (r *Scanner) Scan(
	projectDirectory string,
	moduleName string,
	excludePaths []models.ResolvedPath,
	excludeFileMatchers []*regexp.Regexp,
) ([]models.ProjectFile, error) {
	ctx := resolveContext{
		projectDirectory:    projectDirectory,
		moduleName:          moduleName,
		excludePaths:        excludePaths,
		excludeFileMatchers: excludeFileMatchers,

		tokenSet: token.NewFileSet(),
		results:  []models.ProjectFile{},
	}

	err := r.loadStdPackages()
	if err != nil {
		return nil, fmt.Errorf("failed load std packages info: %w", err)
	}

	err = filepath.Walk(ctx.projectDirectory, func(path string, info os.FileInfo, err error) error {
		return resolveFile(&ctx, path, info, err)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk project tree: %w", err)
	}

	return ctx.results, nil
}

func (r *Scanner) loadStdPackages() error {
	cfg := &packages.Config{
		Mode: packages.NeedName,
	}
	stdList, err := packages.Load(cfg, "std")
	if err != nil {
		return fmt.Errorf("failed load std packages info: %w", err)
	}

	stdPackagesNames = make(map[string]struct{})
	for _, stdPkg := range stdList {
		stdPackagesNames[stdPkg.PkgPath] = struct{}{}
	}

	return nil
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
		return fmt.Errorf("failed to parse go source code at '%s': %w", path, err)
	}

	ctx.results = append(ctx.results, models.ProjectFile{
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
	if _, ok := stdPackagesNames[importPath]; ok {
		return models.ImportTypeStdLib
	}

	if strings.HasPrefix(importPath, ctx.moduleName) {
		return models.ImportTypeProject
	}

	return models.ImportTypeVendor
}
