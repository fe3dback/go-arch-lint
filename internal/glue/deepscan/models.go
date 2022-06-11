package deepscan

import (
	"fmt"
	"go/token"
)

type (
	InjectionMethod struct {
		Name       string // method name (example: `NewProcessor`)
		Definition Source // where method is defined
		Gates      []Gate // method params with interface type
	}

	Source struct {
		Pkg    string   // package name (example: "a")
		Import string   // package full import path (example: "example.com/myProject/internal/a")
		Path   string   // package full abs path (example: "/home/user/go/src/myProject/internal/a")
		Place  Position // exactly place in source code
	}

	Gate struct {
		Name            string           // function param name (func (_a_,b int), name="a")
		Index           int              // function param index (func (a,b bool, c int), for c index=2)
		Definition      Source           // where method param type defined (func (a,b,c _int_))
		Interface       Interface        // used interface for injection
		Implementations []Implementation // all code links to this param
	}

	Interface struct {
		Name       string // interface name
		Definition Source // where interface defined
		GoType     string // interface go type
	}

	Implementation struct {
		Injector Injector // who inject Target to Gate.FunctionName
		Target   Target   // what is injected into Gate.FunctionName
	}

	Injector struct {
		CodeName   string // code expression (not unique)
		Definition Source // where inject occurs
	}

	Target struct {
		StructName string // interface implementation type name
		Definition Source // where this type defined
	}

	Position struct {
		Filename string // filename, if any
		Offset   int    // offset, starting at 0
		Line     int    // line number, starting at 1
		Column   int    // column number, starting at 1 (byte count)
	}
)

func positionFromToken(pos token.Position) Position {
	return Position{
		Filename: pos.Filename,
		Offset:   pos.Offset,
		Line:     pos.Line,
		Column:   pos.Column,
	}
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d", p.Filename, p.Line)
}
