package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const linterVersion = "1.1.0"
const goArchFileSupported = "1"

type versionPayload struct {
	LinterVersion       string
	GoArchFileSupported string
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   cmdNameVersion,
	Short: "Print go arch linter version",
	Run: func(cmd *cobra.Command, args []string) {
		payload := versionPayload{
			LinterVersion:       linterVersion,
			GoArchFileSupported: goArchFileSupported,
		}

		output(payload, func() {
			fmt.Printf("Linter version: %s\n", au.Yellow(payload.LinterVersion))
			fmt.Printf("Supported go arch file versions: %s\n", au.Yellow(payload.GoArchFileSupported))
		})
	},
}
