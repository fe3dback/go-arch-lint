package deepscan

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
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

	// we always expect only one package by path
	parsedPackage := parsedPackages[0]

	// packages.Load reports import/type-checking problems via pkg.Errors,
	// not via returned err. Package with errors has incomplete TypesInfo,
	// so deep scan will silently find nothing. This happens, for example,
	// when local go toolchain writes export data in format, unsupported
	// by golang.org/x/tools version used to build linter.
	if err := packageErrors(parsedPackage); err != nil {
		return nil, err
	}

	ctx.parsedPackages[path] = parsedPackage
	return parsedPackage, nil
}

func packageErrors(pkg *packages.Package) error {
	const maxReported = 5

	if len(pkg.Errors) == 0 {
		return nil
	}

	msgs := make([]string, 0, maxReported+1)
	for i, pkgErr := range pkg.Errors {
		if i >= maxReported {
			msgs = append(msgs, fmt.Sprintf("... and %d more", len(pkg.Errors)-maxReported))
			break
		}

		msgs = append(msgs, pkgErr.Error())
	}

	return fmt.Errorf("package '%s' loaded with errors (deep scan requires fully type-checked code):\n  %s",
		pkg.PkgPath,
		strings.Join(msgs, "\n  "),
	)
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
	mode parser.Mode,
) (files map[string]*ast.File, first error) {
	files = make(map[string]*ast.File)

	parseCtx := parseRecursiveCtx{
		excludedPaths:        excludedPaths,
		excludedFileMatchers: excludedFileMatchers,
		foundFiles:           map[string]struct{}{},
	}

	// WalkDir is used over Walk: it never stats the entries it yields, so an
	// entry that cannot be stat'ed does not abort the walk before the scope
	// check has a chance to skip it.
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		return resolveScopeFile(&parseCtx, path, d, err)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk project tree: %w", err)
	}

	for filePath := range parseCtx.foundFiles {
		fileAst, err := parser.ParseFile(fset, filePath, nil, mode)
		if err != nil {
			return nil, fmt.Errorf("failed parse '%s': %w", filePath, err)
		}

		files[filePath] = fileAst
	}

	return files, nil
}

func resolveScopeFile(ctx *parseRecursiveCtx, path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if _, alreadyExist := ctx.foundFiles[path]; alreadyExist {
		return nil
	}

	if d.IsDir() || !inScope(ctx, path) {
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
