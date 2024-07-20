package container

import (
	"github.com/urfave/cli/v2"

	"github.com/fe3dback/go-arch-lint/v4/internal/app/internal/build"
	"github.com/fe3dback/go-arch-lint/v4/internal/app/internal/flags"
)

func (c *Container) Cli() *cli.App {
	return once(func() *cli.App {
		return &cli.App{
			Name:        "go-arch-lint",
			Usage:       "Golang architecture linter",
			Version:     build.Version,
			Description: "Check all project imports and compare to arch rules defined in yaml file.\nRead full documentation in: https://github.com/fe3dback/go-arch-lint",
			Commands: []*cli.Command{
				c.cliCommandMapping(),
			},
			Flags:    flags.GlobalFlags,
			Compiled: build.CompileTime,
			Authors: []*cli.Author{
				{
					Name:  "fe3dback",
					Email: "fe3dback@pm.me",
				},
			},
			Copyright: "MIT",
		}
	})
}

func (c *Container) cliCommandMapping() *cli.Command {
	return once(func() *cli.Command {
		return &cli.Command{
			Name:        "mapping",
			Aliases:     []string{"ps", "ls", "m"},
			Description: "output mapping table between files and components",
			Action:      c.makeCliCommand("check", c.commandMapping()),
			Flags:       append(flags.GlobalFlags, c.commandMappingFlags()...),
		}
	})
}
