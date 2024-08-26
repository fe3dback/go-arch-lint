package deepscan

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

const parseMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

func cachedPackage(ctx *searchCtx, path string) (*packages.Package, error) {
	if parsedPackage, exist := ctx.parsedPackages[path]; exist {
		return parsedPackage, nil
	}

	cfg := &packages.Config{
		Mode: parseMode,
		Fset: ctx.fileSet,
		Dir:  path,
	}
	parsedPackages, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("failed parse go source: %w", err)
	}

	if len(parsedPackages) == 0 {
		return nil, fmt.Errorf("not found go sources")
	}

	for _, parsedPackage := range parsedPackages {
		// we always expect only one package by path
		ctx.parsedPackages[path] = parsedPackage
		break
	}

	return ctx.parsedPackages[path], nil
}

// isPublicName check that first char in string in uppercase
// so its go public name (like `PublicMethod`)
// return false for `privateMethod`
func astIsPublicName(name string) bool {
	for _, r := range name {
		// check first rune is upper in name
		return unicode.IsUpper(r)
	}

	return false
}

type parseRecursiveCtx struct {
	excludedPaths        []string
	excludedFileMatchers []*regexp.Regexp
	foundFiles           map[string]struct{}
}

func parseRecursive(
	fset *token.FileSet,
	path string,
	excludedPaths []string,
	excludedFileMatchers []*regexp.Regexp,
	filter func(fs.FileInfo) bool,
	mode parser.Mode,
) (pkgs map[string]*ast.Package, first error) {
	pkgs = make(map[string]*ast.Package)

	parseCtx := parseRecursiveCtx{
		excludedPaths:        excludedPaths,
		excludedFileMatchers: excludedFileMatchers,
		foundFiles:           map[string]struct{}{},
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		return resolveScopeFile(&parseCtx, path, info, err)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk project tree: %w", err)
	}

	foundPackages := extractPackagesFromFilePaths(parseCtx.foundFiles)
	for _, analysePackage := range foundPackages {
		found, err := parser.ParseDir(fset, analysePackage, filter, mode)
		if err != nil {
			return nil, fmt.Errorf("failed parse '%s': %w", analysePackage, err)
		}

		for packageID, astPackage := range found {
			pkgs[fmt.Sprintf("%s_%s", analysePackage, packageID)] = astPackage
		}
	}

	return pkgs, nil
}

func extractPackagesFromFilePaths(paths map[string]struct{}) []string {
	r := make(map[string]struct{})

	for path := range paths {
		r[filepath.Dir(path)] = struct{}{}
	}

	list := make([]string, 0, len(r))
	for packagePath := range r {
		list = append(list, packagePath)
	}

	return list
}

func resolveScopeFile(ctx *parseRecursiveCtx, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if _, alreadyExist := ctx.foundFiles[path]; alreadyExist {
		return nil
	}

	if info.IsDir() || !inScope(ctx, path) {
		return nil
	}

	ctx.foundFiles[path] = struct{}{}
	return nil
}

func inScope(ctx *parseRecursiveCtx, path string) bool {
	if filepath.Ext(path) != ".go" {
		return false
	}

	for _, excludePath := range ctx.excludedPaths {
		if strings.HasPrefix(path, excludePath) {
			return false
		}
	}

	for _, matcher := range ctx.excludedFileMatchers {
		if matcher.Match([]byte(path)) {
			return false
		}
	}

	return true
}
