package deepscan

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	goImport       = string
	packageAbsPath = string
	astPackagesMap = map[packageAbsPath]*packages.Package
)

func (s *Searcher) applyImplementations(methods []InjectionMethod) error {
	imports := s.extractImports(methods)
	packagePaths, err := s.findPackagesWithImport(imports)
	if err != nil {
		return fmt.Errorf("failed find go packages in project: %w", err)
	}

	astPackages, err := s.parsePackages(packagePaths)
	if err != nil {
		return fmt.Errorf("failed parse go packages: %w", err)
	}

	s.applyMethodsImplementationsInPackages(methods, astPackages)
	return nil
}

// extract all import path's from all found methods
// next we can search by all source code, when
// *.go files have this imports
func (s *Searcher) extractImports(methods []InjectionMethod) []goImport {
	result := make(map[goImport]struct{}, 0)

	for _, method := range methods {
		result[method.Definition.Import] = struct{}{}
	}

	return mapStrToSlice(result)
}

// fast filter *.go files, who contain any of imports
// and apply file package path to output
func (s *Searcher) findPackagesWithImport(imports []goImport) ([]packageAbsPath, error) {
	err := s.preloadImports()
	if err != nil {
		return nil, fmt.Errorf("failed preload analyse scope imports: %w", err)
	}

	foundPackagesPath := make(map[packageAbsPath]struct{}, 0)
	importsMap := sliceStrToMap(imports)

	for _, astPackage := range s.ctx.parsedImports {
		for filePath, astFile := range astPackage.Files {
			for _, astImportSpec := range astFile.Imports {
				if _, ok := importsMap[strings.Trim(astImportSpec.Path.Value, `"`)]; ok {
					// this file contains useful import
					packagePath := filepath.Dir(filePath)
					foundPackagesPath[packagePath] = struct{}{}
					break
				}
			}
		}

	}

	return mapStrToSlice(foundPackagesPath), nil
}

// parse full ast code and types for every go package provided in paths
func (s *Searcher) parsePackages(paths []packageAbsPath) (astPackagesMap, error) {
	result := make(astPackagesMap)

	for _, packagePath := range paths {
		astPackage, err := cachedPackage(s.ctx, packagePath)
		if err != nil {
			return nil, fmt.Errorf("failed take package '%s': %w", packagePath, err)
		}

		result[packagePath] = astPackage
	}

	return result, nil
}

// find all implementations for each method, and apply it to methods slice items
func (s *Searcher) applyMethodsImplementationsInPackages(methods []InjectionMethod, astPackages astPackagesMap) {
	for _, method := range methods {
		s.applyMethodImplementationsInPackages(&method, astPackages)
	}
}

// find all implementations for method, and apply it
func (s *Searcher) applyMethodImplementationsInPackages(method *InjectionMethod, astPackages astPackagesMap) {
	for _, astPackage := range astPackages {
		for _, astFile := range astPackage.Syntax {
			// default alias is same as package name
			// example `import "path/filepath"`
			// `filepath` - is default alias
			importAlias := path.Base(method.Definition.Import)
			fileHasCalls := false

			// find importAlias for current file
			for _, astImport := range astFile.Imports {
				if strings.Trim(astImport.Path.Value, `"`) != method.Definition.Import {
					continue
				}

				fileHasCalls = true
				if astImport.Name != nil {
					importAlias = astImport.Name.Name
					break
				}
			}

			if !fileHasCalls {
				continue
			}

			s.findFunctionCalls(importAlias, method.Name, astFile, func(callExpr *ast.CallExpr) {
				for gateIndex := range method.Gates {
					gate := &method.Gates[gateIndex]
					callMethod := callExpr.Fun
					callParam := callExpr.Args[gate.Index]

					paramType := astPackage.TypesInfo.TypeOf(callParam)
					targetName, targetPos, valid := s.extractTargetFromCallParam(paramType)
					if !valid {
						// unknown injection type, possible generics or other not known
						// go features on current moment
						// we can extend this function later for new cases
						continue
					}

					targetDefinitions := s.sourceFromToken(targetPos)
					if targetDefinitions.Import == method.Definition.Import {
						// injector use our public interface for typing
						// we exclude this cases from more deep analyse, because of runtime
						// types injection. (we have not known type on compile time)
						continue
					}

					if !targetDefinitions.Place.Valid {
						// invalid target
						// possible is some not importable std const like `errors`
						// or not known ast at this moment
						continue
					}

					gate.Implementations = append(gate.Implementations, Implementation{
						Injector: Injector{
							CodeName:         s.extractCodeFromASTNode(callParam),
							ParamDefinition:  s.sourceFromToken(callParam.Pos()),
							MethodDefinition: s.sourceFromToken(callMethod.Pos()),
						},
						Target: Target{
							StructName: targetName,
							Definition: targetDefinitions,
						},
					})
				}
			})
		}
	}
}

func (s *Searcher) extractCodeFromASTNode(node ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, s.ctx.fileSet, node)
	if err == nil {
		return buf.String()
	}

	return "unknown"
}

func (s *Searcher) extractTargetFromCallParam(t types.Type) (name string, pos token.Pos, valid bool) {
	switch goType := t.(type) {
	case *types.Named:
		return goType.Obj().Name(), goType.Obj().Pos(), true
	case *types.Pointer:
		return s.extractTargetFromCallParam(goType.Elem())
	default:
		return "", pos, false
	}
}

func (s *Searcher) findFunctionCalls(
	packageAlias string,
	functionName string,
	astFile *ast.File,
	onFound func(callExpr *ast.CallExpr),
) {
	ast.Inspect(astFile, func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			// is not function call
			return true
		}

		if callExpr.Fun == nil {
			// not have function in call
			return true
		}

		// check if this outside package function
		// example: `fmt.Println(..)`
		// Println is external function in fmt package
		astFunc, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			// not imported function
			return true
		}

		// x - package name (example: fmt)
		astFuncRef, ok := astFunc.X.(*ast.Ident)
		if !ok {
			return true
		}

		if astFuncRef.Name != packageAlias {
			// not our package
			return true
		}

		if astFunc.Sel.Name != functionName {
			// our package, but another public function
			return true
		}

		// ok, we found caller
		onFound(callExpr)
		return false
	})
}
