package container

import (
	"os"

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
				c.cliCommandCheck(),
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
			ErrWriter: os.Stderr,
		}
	})
}

func (c *Container) cliCommandMapping() *cli.Command {
	return once(func() *cli.Command {
		const name = "mapping"
		return &cli.Command{
			Name:    name,
			Aliases: []string{"ps", "ls", "m"},
			Usage:   "output mapping table between go packages and components",
			Before:  c.makeBeforeCode(),
			Action: c.makeCliCommand(name, func() CommandHandler {
				return c.commandMapping()
			}),
			Flags: append(flags.GlobalFlags, c.commandMappingFlags()...),
		}
	})
}

func (c *Container) cliCommandCheck() *cli.Command {
	return once(func() *cli.Command {
		const name = "check"
		return &cli.Command{
			Name:    name,
			Aliases: []string{"c"},
			Usage:   "compare project *.go files with architectural specification defined in file",
			Before:  c.makeBeforeCode(),
			Action: c.makeCliCommand(name, func() CommandHandler {
				return c.commandCheck()
			}),
			Flags: append(flags.GlobalFlags, c.commandCheckFlags()...),
		}
	})
}
