package cmd

import (
	"github.com/spf13/cobra"
)

const linterVersion = "1.1.0"
const goArchFileSupported = "1"

func cmdVersion(cmd *cobra.Command, _ []string) {
	cmdOutput := payloadVersion{
		LinterVersion:       linterVersion,
		GoArchFileSupported: goArchFileSupported,
	}

	rootFlags := mustFetchFlags(cmd.Context())

	output(
		rootFlags,
		outputPayloadTypeCommandVersion,
		cmdOutput,
		func() {
			versionCmdOutput(rootFlags, cmdOutput)
		},
	)
}
