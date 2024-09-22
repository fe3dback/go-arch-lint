package view

import (
	_ "embed"
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/commands/check"
	"github.com/fe3dback/go-arch-lint-sdk/commands/mapping"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

//go:embed view_error.gohtml
var viewError []byte

//go:embed view_check.gohtml
var viewCheck []byte

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

var Templates = map[string][]byte{
	tpl(models.CmdStdoutErrorOut{}): viewError,
	tpl(mapping.Out{}):              viewMapping,
	tpl(check.Out{}):                viewCheck,
	//tpl(models.CmdGraphOut{}):       viewGraph,
	//tpl(models.CmdSchemaOut{}):      viewSchema,
	//tpl(models.CmdSelfInspectOut{}): viewSelfInspect,
	//tpl(models.CmdVersionOut{}):     viewVersion,
}

func tpl(model interface{}) string {
	return fmt.Sprintf("%T", model)
}
