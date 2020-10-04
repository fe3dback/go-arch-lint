package cmd

import "fmt"

func versionCmdOutput(flags *rootInput, cmdOutput payloadVersion) {
	au := flags.au

	fmt.Printf("Linter version: %s\n", au.Yellow(cmdOutput.LinterVersion))
	fmt.Printf("Supported go arch file versions: %s\n", au.Yellow(cmdOutput.GoArchFileSupported))
}
