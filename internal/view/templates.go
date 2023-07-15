package view

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

var Templates = map[string]string{
	tpl(models.CmdErrorOut{}):       Error,
	tpl(models.CmdVersionOut{}):     Version,
	tpl(models.CmdSelfInspectOut{}): SelfInspect,
	tpl(models.CmdCheckOut{}):       Check,
	tpl(models.CmdMappingOut{}):     Mapping,
	tpl(models.CmdSchemaOut{}):      Schema,
	tpl(models.CmdGraphOut{}):       Graph,
}

func tpl(model interface{}) string {
	return fmt.Sprintf("%T", model)
}
