package container

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/commands/check"
	"github.com/fe3dback/go-arch-lint/v4/internal/commands/mapping"
)

func (c *Container) commandMapping() *mapping.Command {
	return once(func() *mapping.Command {
		return mapping.NewCommand(
			c.sdk(),
			c.serviceSpecFetcher(),
		)
	})
}

func (c *Container) commandMappingFlags() []cli.Flag {
	return once(func() []cli.Flag {
		return mapping.Flags
	})
}

func (c *Container) commandCheck() *check.Command {
	return once(func() *check.Command {
		return check.NewCommand(
			c.sdk(),
			c.serviceSpecFetcher(),
		)
	})
}

func (c *Container) commandCheckFlags() []cli.Flag {
	return once(func() []cli.Flag {
		return check.Flags
	})
}
