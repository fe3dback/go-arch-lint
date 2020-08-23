package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fe3dback/go-arch-lint/checker"
	"github.com/fe3dback/go-arch-lint/files"
	"github.com/fe3dback/go-arch-lint/spec"
	"github.com/logrusorgru/aurora/v3"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}
	}()

	flags, err := flagsParse()
	if err != nil {
		panic(err)
	}

	arch, err := spec.NewArch(
		flags.pathArchFile,
		flags.moduleName,
		flags.pathProjectDirectory,
	)
	if err != nil {
		panic(err)
	}

	resolver := createResolver(arch, flags)
	projectFiles, err := resolver.Resolve()
	if err != nil {
		panic(err)
	}

	archChecker := checker.NewChecker(
		flags.pathProjectDirectory,
		arch,
		projectFiles,
	)

	result := archChecker.Check()
	au := aurora.NewAurora(flags.useColorsOutput)

	if result.IsOk() {
		fmt.Printf("%s\n", au.Green("OK - No warnings found"))
		return
	}

	writeWarnings(
		result,
		os.Stderr,
		flags.maxWarningsOutput,
		flags.useColorsOutput,
	)

	panic(fmt.Sprintf("\nWarnings found: %d", au.Yellow(result.TotalCount())))
}

func writeWarnings(res checker.CheckResult, writer io.Writer, maxWarnings int, useColors bool) {
	depsWarnings := res.DependencyWarnings()
	matchWarnings := res.NotMatchedWarnings()

	outputCount := 0
	out := func(warn string, force bool) {
		if (outputCount >= maxWarnings) && !force {
			return
		}

		_, err := writer.Write([]byte(warn))
		if err != nil {
			panic(fmt.Sprintf("can`t write warnings to writer interface: %v", err))
		}

		if !force {
			outputCount++
		}
	}

	au := aurora.NewAurora(useColors)

	for _, warning := range depsWarnings {
		out(fmt.Sprintf("[WARN] Component '%s': file '%s' shouldn't depend on '%s'\n",
			au.Green(warning.ComponentName),
			au.Cyan(warning.FileRelativePath),
			au.Yellow(warning.ResolvedImportName),
		), false)
	}

	for _, warning := range matchWarnings {
		out(fmt.Sprintf("[WARN] File '%s' not attached to any component in archfile\n",
			au.Cyan(warning.FileRelativePath),
		), false)
	}

	truncatedWarnings := res.TotalCount() - outputCount
	if truncatedWarnings >= 1 {
		out(fmt.Sprintf("%d warning truncated..\n", truncatedWarnings), true)
	}
}

func createResolver(arch *spec.Arch, flags flags) *files.Resolver {
	excludePaths := make([]string, 0)
	for _, excludeDir := range arch.Exclude {
		excludePaths = append(excludePaths, excludeDir.AbsPath)
	}

	resolver := files.NewResolver(
		flags.pathProjectDirectory,
		flags.moduleName,
		excludePaths,
		arch.ExcludeFilesMatcher,
	)

	return resolver
}
