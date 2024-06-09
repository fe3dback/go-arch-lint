package deepscan

import (
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type (
	InjectionMethod struct {
		Name       string // method name (example: `NewProcessor`)
		Definition Source // where method is defined
		Gates      []Gate // method params with interface type
	}

	Gate struct {
		MethodName         string           // function name (func Hello(a,b int), name="Hello")
		ParamName          string           // function param name (func (_a_,b int), name="a")
		Index              int              // function param index (func (a,b bool, c int), for c index=2)
		MethodDefinition   Source           // where method is defined
		ArgumentDefinition Source           // where method param type defined (func (a,b,c _int_))
		Interface          Interface        // used interface for injection
		Implementations    []Implementation // all code links to this param
		IsVariadic         bool             // function param is variadic (func (a bool, nums ...int))
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
		CodeName         string // code expression (not unique)
		MethodDefinition Source // where method is called
		ParamDefinition  Source // where param is passed to method (injection occurs)
	}

	Target struct {
		StructName string // interface implementation type name
		Definition Source // where this type defined
	}

	Source struct {
		Pkg    string           // package name (example: "a")
		Import string           // package full import path (example: "example.com/myProject/internal/a")
		Path   string           // package full abs path (example: "/home/user/go/src/myProject/internal/a")
		Place  common.Reference // exactly place in source code
	}
)
