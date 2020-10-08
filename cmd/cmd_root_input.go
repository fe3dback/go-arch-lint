package cmd

import (
	"fmt"
	"strings"

	"github.com/fe3dback/go-arch-lint/models"

	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

// flag names
const (
	flagNameOutputColors      = "output-color"
	flagNameOutputJsonOneLine = "output-json-one-line"
	flagNameOutputType        = "output-type"
	flagNameMaxWarnings       = "max-warnings"
	flagNameProjectPath       = "project-path"
	flagNameArchFile          = "arch-file"
	flagNameAliasJson         = "json"
)

// output types
const (
	outputTypeDefault outputType = "default"
	outputTypeASCII   outputType = "ascii"
	outputTypeJSON    outputType = "json"
)

const defaultArchFileName = ".go-arch-lint.yml"

type (
	rootInput struct {
		au                   aurora.Aurora
		useColors            bool
		outputType           outputType
		outputJsonOneLine    bool
		maxWarnings          int
		projectRootDirectory string
		archFile             string
		useJsonAlias         bool
	}
)

func newDefaultFlags() *rootInput {
	return &rootInput{
		au:                   nil,
		useColors:            true,
		useJsonAlias:         false,
		outputType:           outputTypeDefault,
		outputJsonOneLine:    false,
		maxWarnings:          512,
		projectRootDirectory: "",
		archFile:             defaultArchFileName,
	}
}

func parseFlags(rootCmd *cobra.Command, input *rootInput) *rootInput {
	// color-output
	rootCmd.PersistentFlags().BoolVar(
		&input.useColors,
		flagNameOutputColors,
		input.useColors,
		"use ANSI colors in terminal output",
	)

	// color-output
	rootCmd.PersistentFlags().BoolVar(
		&input.outputJsonOneLine,
		flagNameOutputJsonOneLine,
		input.outputJsonOneLine,
		"format JSON as single line payload (without line breaks), only for json output type",
	)

	// output-type
	rootCmd.PersistentFlags().StringVar(
		&input.outputType,
		flagNameOutputType,
		input.outputType,
		fmt.Sprintf("type of command output, variants: [%s]", strings.Join(outputTypeVariantsConst, ", ")),
	)

	// json alias
	rootCmd.PersistentFlags().BoolVar(
		&input.useJsonAlias,
		flagNameAliasJson,
		input.useJsonAlias,
		fmt.Sprintf("(alias for --%s=%s)",
			flagNameOutputType,
			outputTypeJSON,
		),
	)

	// max warnings
	rootCmd.PersistentFlags().IntVar(
		&input.maxWarnings,
		flagNameMaxWarnings,
		input.maxWarnings,
		"max number of warnings to output",
	)

	// project root directory
	rootCmd.PersistentFlags().StringVar(
		&input.projectRootDirectory,
		flagNameProjectPath,
		input.projectRootDirectory,
		fmt.Sprintf("absolute path to project directory (where '%s' is located)", defaultArchFileName),
	)

	// go arch file
	rootCmd.PersistentFlags().StringVar(
		&input.archFile,
		flagNameArchFile,
		input.archFile,
		"arch file path",
	)

	// init
	cobra.OnInitialize(func() {
		input.au = aurora.NewAurora(input.useColors)

		processAliases(input)
		assertFlagOutputTypeValid(input)
		assertFlagMaxWarningsValid(input)
	})

	return input
}

func processAliases(input *rootInput) {
	processAliasJson(input)
}

func processAliasJson(input *rootInput) {
	if !input.useJsonAlias {
		return
	}

	if input.outputType != outputTypeDefault && input.outputType != outputTypeJSON {
		panic(models.NewUserSpaceError(
			fmt.Sprintf("flag --%s not compatible with --%s",
				flagNameAliasJson,
				flagNameOutputType,
			),
		))
	}

	input.outputType = outputTypeJSON
}

func assertFlagOutputTypeValid(input *rootInput) {
	if input.outputType == outputTypeDefault {
		input.outputType = outputTypeASCII
	}

	for _, variant := range outputTypeVariantsConst {
		if input.outputType == variant {
			return
		}
	}

	panic(models.NewUserSpaceError(fmt.Sprintf("unknown output-type: %s", input.outputType)))
}

func assertFlagMaxWarningsValid(input *rootInput) {
	const warningsRangeMin = 1
	const warningsRangeMax = 32768

	if input.maxWarnings < warningsRangeMin || input.maxWarnings > warningsRangeMax {
		panic(models.NewUserSpaceError(
			fmt.Sprintf("flag '%s' should by in range [%d .. %d]",
				flagNameMaxWarnings,
				warningsRangeMin,
				warningsRangeMax,
			),
		))
	}
}
