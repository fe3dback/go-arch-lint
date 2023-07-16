package deepscan

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func (s *Searcher) extractMethodsFromPackage(astPackage *packages.Package) ([]InjectionMethod, error) {
	result := make([]InjectionMethod, 0)

	for _, astFile := range astPackage.Syntax {
		methods, err := s.extractMethodsFromFile(astPackage, astFile)
		if err != nil {
			return nil, fmt.Errorf("failed extract public methods from '%s': %w", astFile.Name.String(), err)
		}

		result = append(result, methods...)
	}

	return result, nil
}

func (s *Searcher) extractMethodsFromFile(astPackage *packages.Package, astFile *ast.File) ([]InjectionMethod, error) {
	list := make([]InjectionMethod, 0)

	for _, iDecl := range astFile.Decls {
		decl, ok := iDecl.(*ast.FuncDecl)
		if !ok {
			// find only AST go methods, example: `func a()`
			continue
		}

		if !astIsPublicName(decl.Name.Name) {
			// exclude private methods
			continue
		}

		gates := s.extractMethodGates(astPackage, decl)
		if len(gates) == 0 {
			// this method not have interface params (gates)
			// so nothing can be injected into
			// and no reason for analyse it
			continue
		}

		list = append(list, InjectionMethod{
			Name:       decl.Name.Name,
			Definition: s.sourceFromToken(decl.Name.Pos()),
			Gates:      gates,
		})
	}

	return list, nil
}

func (s *Searcher) extractMethodGates(astPackage *packages.Package, method *ast.FuncDecl) []Gate {
	fields := method.Type.Params.List
	params := make([]Gate, 0, len(fields))
	typeIndex := -1

	for _, field := range fields {
		paramType := astPackage.TypesInfo.TypeOf(field.Type)
		for _, fieldIdent := range field.Names {
			typeIndex++

			if fieldIdent.Name == "_" {
				// argument not used, so we not use any
				// implementation logic, and not depend on it
				continue
			}

			interfaceName, pos, isInterface := s.extractInterfaceName(paramType)
			if !isInterface {
				continue
			}

			if !pos.IsValid() {
				// invalid pos, its anonymous `interface{}`
				// or some kind of this staff
				pos = field.Pos()
			}

			params = append(params, Gate{
				MethodName:         method.Name.Name,
				ParamName:          fieldIdent.Name,
				Index:              typeIndex,
				MethodDefinition:   s.sourceFromToken(method.Pos()),
				ArgumentDefinition: s.sourceFromToken(field.Pos()),
				Interface: Interface{
					Name:       interfaceName,
					Definition: s.sourceFromToken(pos),
					GoType:     paramType.String(),
				},
			})
		}
	}

	return params
}

func (s *Searcher) extractInterfaceName(t types.Type) (name string, ref token.Pos, isInterface bool) {
	switch goType := t.(type) {
	// anon interfaces: `func(a interface{})`
	case *types.Interface:
		return t.String(), ref, true

	// named type or interface: `func(a myInterface)` or `func (a int)`
	case *types.Named:
		if _, _, isInterface := s.extractInterfaceName(goType.Underlying()); !isInterface {
			return "", ref, false
		}

		return goType.Obj().Name(), goType.Obj().Pos(), true

	// pointer to type: `func(a *int)`, possible can point to interface
	// but is useless in real code, so always skip this params
	case *types.Pointer:
		return "", ref, false

	case *types.Map:
		return s.extractInterfaceName(goType.Elem())

	// `func(a []int)` or `func(a []myInterface)`
	// can be slice of interface, so need to check it
	case *types.Slice:
		return s.extractInterfaceName(goType.Elem())

	// `func(a [5]int)` or `func(a [5]myInterface)`
	// can be array of interface, so need to check it
	case *types.Array:
		return s.extractInterfaceName(goType.Elem())

	// `func(a chan int)` or `func(a chan myInterface)`
	// need to check operand, and where is interface
	// we have 3 possible chan options
	// <- chan :: its read only channel (implementation can be injected)
	// chan <- :: write only chanel (nothing can be injected)
	// chan    :: r/w, possible can inject
	// if chan cannot inject anything, we skip it from analyse
	case *types.Chan:
		deepName, deepPos, isInterface := s.extractInterfaceName(goType.Elem())
		if !isInterface {
			return "", ref, false
		}

		if goType.Dir() == types.SendOnly {
			// nothing be injected into write only interface
			return "", ref, false
		}

		return deepName, deepPos, true

	// not interface
	default:
		return "", ref, false
	}
}
