package container

import (
	"github.com/fe3dback/go-arch-lint/internal/operations/version"
	"github.com/spf13/cobra"
)

func (c *Container) commandVersion() (*cobra.Command, runner) {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print go arch linter version",
		Long:  "show version, build time and commit hash of current build",
	}

	return cmd, func(_ *cobra.Command) (any, error) {
		return c.commandVersionOperation().Behave()
	}
}

func (c *Container) commandVersionOperation() *version.Operation {
	return version.NewOperation(
		c.version,
		c.buildTime,
		c.commitHash,
	)
}
