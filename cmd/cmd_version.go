package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const linterVersion = "1.1.0"
const goArchFileSupported = "1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   cmdNameVersion,
	Short: "Print go arch linter version",
	Run: func(cmd *cobra.Command, args []string) {
		payload := payloadVersion{
			LinterVersion:       linterVersion,
			GoArchFileSupported: goArchFileSupported,
		}

		output(outputPayloadTypeCommandVersion, payload, func() {
			fmt.Printf("Linter version: %s\n", au.Yellow(payload.LinterVersion))
			fmt.Printf("Supported go arch file versions: %s\n", au.Yellow(payload.GoArchFileSupported))
		})
	},
}
