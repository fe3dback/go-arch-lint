package schema

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

const (
	flagVersion = "version"
)

type (
	processorFn = func(schema models.FlagsSchema) error

	CommandAssembler struct {
		jsonSchemaProvider JsonSchemaProvider
		processorFn        processorFn
		localFlags         *localFlags
	}

	localFlags struct {
		Version int
	}
)

func NewSchemaCommandAssembler(
	jsonSchemaProvider JsonSchemaProvider,
	processorFn processorFn,
) *CommandAssembler {
	return &CommandAssembler{
		jsonSchemaProvider: jsonSchemaProvider,
		processorFn:        processorFn,
		localFlags: &localFlags{
			Version: 0,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema",
		Short:   "json schema for arch file inspection",
		Long:    "useful for integrations with ide's and editor plugins",
		PreRunE: c.prePersist,
		RunE:    c.invoke,
	}

	c.assembleFlags(cmd)

	return cmd
}

func (c *CommandAssembler) invoke(_ *cobra.Command, _ []string) error {
	input, err := c.assembleInput()
	if err != nil {
		return fmt.Errorf("failed to assemble input params: %w", err)
	}

	return c.processorFn(input)
}

func (c *CommandAssembler) prePersist(cmd *cobra.Command, _ []string) error {
	version, err := cmd.Flags().GetInt(flagVersion)
	if err != nil {
		return failedToGetFlag(err, flagVersion)
	}

	c.localFlags.Version = version

	return nil
}

func (c *CommandAssembler) assembleInput() (models.FlagsSchema, error) {
	jsonSchema, err := c.jsonSchemaProvider.Provide(c.localFlags.Version)
	if err != nil {
		return models.FlagsSchema{}, fmt.Errorf("failed to provide json schema: %w", err)
	}

	return models.FlagsSchema{
		Version:    c.localFlags.Version,
		JsonSchema: jsonSchema,
	}, nil
}
