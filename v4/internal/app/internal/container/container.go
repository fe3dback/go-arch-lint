package container

import "github.com/urfave/cli/v2"

type Container struct {
	cCtx *cli.Context
}

func NewContainer() *Container {
	return &Container{}
}
