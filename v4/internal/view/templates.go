package view

import (
	_ "embed"
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/mapping"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

//go:embed view_error.gohtml
var viewError []byte

////go:embed view_check.gohtml
//var viewCheck []byte
//
////go:embed view_error.gohtml
//var viewError []byte
//
////go:embed view_graph.gohtml
//var viewGraph []byte

//go:embed view_mapping.gohtml
var viewMapping []byte

////go:embed view_schema.gohtml
//var viewSchema []byte
//
////go:embed view_self_inspect.gohtml
//var viewSelfInspect []byte
//
////go:embed view_version.gohtml
//var viewVersion []byte

var Templates = map[string]string{
	tpl(models.CmdStdoutErrorOut{}): string(viewError),
	//tpl(models.CmdCheckOut{}):       string(viewCheck),
	//tpl(models.CmdErrorOut{}):       string(viewError),
	//tpl(models.CmdGraphOut{}):       string(viewGraph),
	tpl(mapping.Out{}): string(viewMapping),
	//tpl(models.CmdSchemaOut{}):      string(viewSchema),
	//tpl(models.CmdSelfInspectOut{}): string(viewSelfInspect),
	//tpl(models.CmdVersionOut{}):     string(viewVersion),
}

func tpl(model interface{}) string {
	return fmt.Sprintf("%T", model)
}
