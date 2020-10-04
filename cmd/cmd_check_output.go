package cmd

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/models"
)

func checkCmdOutputAscii(flags *rootInput, input checkCmdInput, output payloadTypeCommandCheck) {
	au := flags.au

	fmt.Printf("module: %s\n", au.Green(input.settingsModuleName))

	if output.ExecutionError != "" {
		for _, warning := range output.ExecutionWarnings {
			if src := warning.SourceCode; src != nil {
				fmt.Printf("[Archfile] %s:\n%s\n",
					au.Yellow(warning.Text),
					warning.SourceCode,
				)

				continue
			}

			fmt.Printf("[Archfile] %s\n", au.Yellow(warning.Text))
		}

		panic(models.NewUserSpaceError(output.ExecutionError))
	}

	if !output.ArchHasWarnings {
		fmt.Printf("%s\n", au.Green("OK - No warnings found"))

		return
	}

	outputCount := 0
	for _, warning := range output.ArchWarningsDeps {
		fmt.Printf("[WARN] Component '%s': file '%s' shouldn't depend on '%s'\n",
			au.Green(warning.ComponentName),
			au.Cyan(warning.FileRelativePath),
			au.Yellow(warning.ResolvedImportName),
		)
		outputCount++
	}

	for _, warning := range output.ArchWarningsNotMatched {
		fmt.Printf("[WARN] File '%s' not attached to any component in archfile\n",
			au.Cyan(warning.FileRelativePath),
		)
		outputCount++
	}

	truncatedWarnings := len(output.ArchWarningsDeps) - outputCount
	if truncatedWarnings >= 1 {
		fmt.Printf("%d warning truncated..\n", truncatedWarnings)
	}

	fmt.Println() // write empty line
	panic(models.NewUserSpaceError(fmt.Sprintf("warnings found: %d", au.Yellow(len(output.ArchWarningsDeps)))))
}
