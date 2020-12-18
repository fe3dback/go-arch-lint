package version

import (
	"github.com/spf13/cobra"
)

type (
	processorFn = func() error

	CommandAssembler struct {
		processorFn processorFn
	}
)

func NewVersionCommandAssembler(processorFn processorFn) *CommandAssembler {
	return &CommandAssembler{
		processorFn: processorFn,
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print go arch linter version",
		Long:  "show version, build time and commit hash of current build",
		RunE:  c.invoke,
	}

	return cmd
}

func (c *CommandAssembler) invoke(_ *cobra.Command, _ []string) error {
	return c.processorFn()
}
