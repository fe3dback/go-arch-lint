package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/checker"
	"github.com/fe3dback/go-arch-lint/files"
	"github.com/fe3dback/go-arch-lint/spec"
	"github.com/spf13/cobra"
)

var (
	settingsProjectDirectory string
	settingsGoArchFilePath   string
	settingsGoModFilePath    string
	settingsModuleName       string
)

type (
	checkPayload struct {
		HasWarnings        bool
		WarningsDeps       []checker.WarningDep
		WarningsNotMatched []checker.WarningNotMatched
	}
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   cmdNameCheck,
	Short: "check project architecture by yaml file",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkCmdAssertFlagProjectPathValid()
		checkCmdAssertFlagGoModuleValid()
	},
	Run: func(cmd *cobra.Command, args []string) {
		payload := checkCmdArch()
		output(payload, func() {
			checkCmdOutputAscii(payload)
		})
	},
}

func checkCmdOutputAscii(payload checkPayload) {
	fmt.Printf("used arch file: %s\n", au.Green(settingsGoArchFilePath))
	fmt.Printf("        module: %s\n", au.Green(settingsModuleName))

	if !payload.HasWarnings {
		fmt.Printf("%s\n", au.Green("OK - No warnings found"))

		return
	}

	outputCount := 0
	for _, warning := range payload.WarningsDeps {
		fmt.Printf("[WARN] Component '%s': file '%s' shouldn't depend on '%s'\n",
			au.Green(warning.ComponentName),
			au.Cyan(warning.FileRelativePath),
			au.Yellow(warning.ResolvedImportName),
		)
		outputCount++
	}

	for _, warning := range payload.WarningsNotMatched {
		fmt.Printf("[WARN] File '%s' not attached to any component in archfile\n",
			au.Cyan(warning.FileRelativePath),
		)
		outputCount++
	}

	truncatedWarnings := len(payload.WarningsDeps) - outputCount
	if truncatedWarnings >= 1 {
		fmt.Printf("%d warning truncated..\n", truncatedWarnings)
	}

	fmt.Println()
	halt(fmt.Errorf("warnings found: %d", au.Yellow(len(payload.WarningsDeps))))
}

func checkCmdAssertFlagProjectPathValid() {
	settingsProjectDirectory = flagProjectPath
	if settingsProjectDirectory == "" {
		panic(fmt.Sprintf("flag '%s' should by set", flagNameProjectPath))
	}

	settingsProjectDirectory = filepath.Clean(settingsProjectDirectory)

	// check arch file
	settingsGoArchFilePath = filepath.Clean(fmt.Sprintf("%s/%s", settingsProjectDirectory, archFileName))
	_, err := os.Stat(settingsGoArchFilePath)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("not found archfile in '%s'", settingsGoArchFilePath))
	}

	// check go.mod
	settingsGoModFilePath = filepath.Clean(fmt.Sprintf("%s/%s", settingsProjectDirectory, goModFileName))
	_, err = os.Stat(settingsGoModFilePath)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("not found project '%s' in '%s'", goModFileName, settingsGoModFilePath))
	}
}

func checkCmdAssertFlagGoModuleValid() {
	moduleName, err := getModuleNameFromGoModFile(settingsGoModFilePath)
	if err != nil {
		panic(fmt.Sprintf("failed get module name: %s", err))
	}

	settingsModuleName = moduleName
}

func checkCmdArch() checkPayload {
	arch, err := spec.NewArch(
		settingsGoArchFilePath,
		settingsModuleName,
		settingsProjectDirectory,
	)
	if err != nil {
		panic(err)
	}

	resolver := checkCmdCreateResolver(arch)
	projectFiles, err := resolver.Resolve()
	if err != nil {
		panic(err)
	}

	archChecker := checker.NewChecker(
		settingsProjectDirectory,
		arch,
		projectFiles,
	)

	payload := checkPayload{
		HasWarnings:        false,
		WarningsDeps:       []checker.WarningDep{},
		WarningsNotMatched: []checker.WarningNotMatched{},
	}

	result := archChecker.Check()
	if result.IsOk() {
		return payload
	}

	checkCmdWriteWarnings(
		result,
		&payload,
		flagMaxWarnings,
	)

	return payload

}

func checkCmdWriteWarnings(res checker.CheckResult, payload *checkPayload, maxWarnings int) {
	outputCount := 0

	// append deps
	for _, dep := range res.DependencyWarnings() {
		if outputCount >= maxWarnings {
			break
		}

		payload.WarningsDeps = append(payload.WarningsDeps, dep)
		outputCount++
	}

	// append not matched
	for _, notMatched := range res.NotMatchedWarnings() {
		if outputCount >= maxWarnings {
			break
		}

		payload.WarningsNotMatched = append(payload.WarningsNotMatched, notMatched)
		outputCount++
	}

	if outputCount > 0 {
		payload.HasWarnings = true
	}
}

func checkCmdCreateResolver(arch *spec.Arch) *files.Resolver {
	excludePaths := make([]string, 0)
	for _, excludeDir := range arch.Exclude {
		excludePaths = append(excludePaths, excludeDir.AbsPath)
	}

	resolver := files.NewResolver(
		settingsProjectDirectory,
		settingsModuleName,
		excludePaths,
		arch.ExcludeFilesMatcher,
	)

	return resolver
}