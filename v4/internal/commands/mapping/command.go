package mapping

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Command struct {
	operation operation
}

func NewCommand(operation operation) *Command {
	return &Command{
		operation: operation,
	}
}

func (c *Command) Execute(cCtx *cli.Context) (any, error) {
	in := c.parseIn(cCtx)

	// todo: validation

	out, err := c.operation.Mapping(in)
	if err != nil {
		return "", err
	}

	// todo: map errors (?)

	return out, nil
}

func (c *Command) parseIn(cCtx *cli.Context) models.CmdMappingIn {
	in := models.CmdMappingIn{}
	in.ProjectPath = cCtx.Path(flagProjectPath)
	in.ArchFile = cCtx.Path(flagArchConfigRelativePath)
	in.Scheme = cCtx.String(flagScheme)

	return in
}
