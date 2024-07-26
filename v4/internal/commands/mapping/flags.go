package mapping

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/pkg/xflags"
)

const (
	flagScheme = "scheme"
)

var Flags = []cli.Flag{
	xflags.CreateEnumFlag(
		flagScheme,
		[]string{"s"},
		"display scheme",
		models.CmdMappingSchemesValues,
		models.CmdMappingSchemeList,
		models.FlagCategoryCommand,
	),
}
