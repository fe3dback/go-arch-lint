package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

func (c *CommandAssembler) prepareOutputType(cmd *cobra.Command) (models.OutputType, error) {
	outputType, err := cmd.Flags().GetString(flagOutputType)
	if err != nil {
		return "", failedToGetFlag(err, flagOutputType)
	}

	useJsonAlias, err := cmd.Flags().GetBool(flagAliasJson)
	if err != nil {
		return "", failedToGetFlag(err, flagAliasJson)
	}

	// alias preprocessor
	if useJsonAlias {
		if outputType != models.OutputTypeDefault && outputType != models.OutputTypeJSON {
			return "", fmt.Errorf("flag --%s not compatible with --%s=%s",
				flagAliasJson,
				flagOutputType,
				outputType,
			)
		}

		outputType = models.OutputTypeJSON
	}

	// fallback to default's
	if outputType == models.OutputTypeDefault {
		outputType = models.OutputTypeASCII
	}

	// validate
	for _, variant := range models.OutputTypeVariantsConst {
		if outputType == variant {
			return outputType, nil
		}
	}

	return "", fmt.Errorf("unknown output-type: %s", outputType)
}
