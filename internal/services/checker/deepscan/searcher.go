package deepscan

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

type (
	// abs path to directory
	absPath = string

	// parsed packages
	packageCache = map[absPath]*packages.Package
)

type (
	Searcher struct {
		ctx *searchCtx

		mux sync.Mutex
	}

	searchCtx struct {
		// current search ctx
		criteria Criteria

		// hold all parsed packages in analyse scope
		// but only with imports declarations
		// used only for fast filter possible params
		parsedImports []*ast.Package

		// hold all already parsed ast packages
		// this holds all package meta information
		// like AST, types, etc...
		parsedPackages packageCache

		// parsed fileset
		fileSet *token.FileSet
	}
)

func NewSearcher() *Searcher {
	return &Searcher{
		ctx: &searchCtx{
			parsedImports:  []*ast.Package{},
			parsedPackages: map[absPath]*packages.Package{},
			fileSet:        token.NewFileSet(),
		},
	}
}

// Usages share same packages cache for every function call
// so it`s good idea to check every package in project
// with same Searcher instance
//
// This method will find all package functions with interfaces
// and link it to all callers, with implementations
// it will skip:
//   - methods without interface (not injectable)
//   - private methods (nobody outside can call it)
//   - only write chan (func (ch chan<-) (our code send something, so we not depend on implementations)
//   - with placeholder param names (func (_ myInterface)), nobody can use _, so code not depend on interface
//
// Can`t search from multiple goroutines, but safe for concurrent use (mutex inside)
func (s *Searcher) Usages(c Criteria) ([]InjectionMethod, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// current ctx
	s.ctx.criteria = c

	// search
	astPackage, err := cachedPackage(s.ctx, c.packagePath)
	if err != nil {
		return nil, fmt.Errorf("failed get package at '%s': %w", c.packagePath, err)
	}

	methods, err := s.extractMethodsFromPackage(astPackage)
	if err != nil {
		return nil, fmt.Errorf("failed extract methods from package at '%s': %w", c.packagePath, err)
	}

	err = s.applyImplementations(methods)
	if err != nil {
		return nil, fmt.Errorf("failed apply implementations info for found methods: %w", err)
	}

	return methods, nil
}

func (s *Searcher) sourceFromToken(pos token.Pos) Source {
	place := positionFromToken(s.ctx.fileSet.Position(pos))
	absPath := filepath.Dir(place.File)
	importRef := s.pathToImport(absPath)
	pkg := path.Base(importRef)

	return Source{
		Pkg:    pkg,
		Import: importRef,
		Path:   absPath,
		Place:  place,
	}
}

func (s *Searcher) pathToImport(packagePath string) string {
	packagePath = strings.TrimPrefix(packagePath, s.ctx.criteria.moduleRootPath)
	packagePath = strings.TrimPrefix(packagePath, string(filepath.Separator))
	packagePath = strings.ReplaceAll(packagePath, string(filepath.Separator), "/")

	return fmt.Sprintf("%s/%s", s.ctx.criteria.moduleName, packagePath)
}

func (s *Searcher) preloadImports() error {
	if len(s.ctx.parsedImports) != 0 {
		return nil
	}

	found, err := parseRecursive(
		s.ctx.fileSet,
		s.ctx.criteria.analyseScope,
		s.ctx.criteria.excludePaths,
		s.ctx.criteria.excludeFileMatchers,
		nil,
		parser.ImportsOnly,
	)
	if err != nil {
		return fmt.Errorf("failed parse imports in scope '%s': %w", s.ctx.criteria.analyseScope, err)
	}

	for _, scopePackage := range found {
		s.ctx.parsedImports = append(s.ctx.parsedImports, scopePackage)
	}

	return nil
}
