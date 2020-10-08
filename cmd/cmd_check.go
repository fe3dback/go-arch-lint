package cmd

import (
	"github.com/spf13/cobra"
)

func cmdCheck(cmd *cobra.Command, _ []string) {
	rootFlags := mustFetchFlags(cmd.Context())
	cmdInput := checkCmdAssembleCommandInput(rootFlags)
	cmdOutput := checkCmdSortOutput(
		checkCmdProcess(rootFlags, cmdInput),
	)

	output(
		cmdOutput.ExecutionError != "" || cmdOutput.ArchHasWarnings,
		rootFlags,
		outputPayloadTypeCommandCheck,
		cmdOutput,
		func() {
			checkCmdOutputAscii(rootFlags, cmdInput, cmdOutput)
		},
	)
}
