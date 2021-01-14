package root

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"

	"github.com/spf13/cobra"
)

func (c *CommandAssembler) prepareOutputType(cmd *cobra.Command) (models.OutputType, error) {
	outputType, err := cmd.Flags().GetString(flagOutputType)
	if err != nil {
		return "", failedToGetFlag(err, flagOutputType)
	}

	useJSONAlias, err := cmd.Flags().GetBool(flagAliasJSON)
	if err != nil {
		return "", failedToGetFlag(err, flagAliasJSON)
	}

	// alias preprocessor
	if useJSONAlias {
		if outputType != models.OutputTypeDefault && outputType != models.OutputTypeJSON {
			return "", fmt.Errorf("flag --%s not compatible with --%s=%s",
				flagAliasJSON,
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
