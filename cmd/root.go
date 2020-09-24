package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

const (
	cmdNameCheck   = "check"
	cmdNameVersion = "version"
)

var rootCmd = &cobra.Command{
	Use:   "go-arch-lint",
	Short: "Golang architecture linter",
	Long: `
	Check all project imports and compare to arch rules defined in yaml file.
	Read full documentation in: https://github.com/fe3dback/go-arch-lint`,
	Run: func(cmd *cobra.Command, args []string) {},
}

// flag names
const (
	flagNameOutputColors      = "output-color"
	flagNameOutputJsonOneLine = "output-json-one-line"
	flagNameOutputType        = "output-type"
	flagNameMaxWarnings       = "max-warnings"
	flagNameProjectPath       = "project-path"
	flagNameAliasJson         = "json"
)

var (
	// flags
	flagUseColors         bool
	flagOutputType        outputType
	flagOutputJsonOneLine bool
	flagMaxWarnings       int
	flagProjectPath       string
	flagAliasJson         bool

	// global services
	au aurora.Aurora
)

var (
	outputTypeVariants = []string{outputTypeASCII, outputTypeJSON}
)

type (
	outputType = string
)

const (
	outputTypeDefault outputType = "default"
	outputTypeASCII   outputType = "ascii"
	outputTypeJSON    outputType = "json"
)

const (
	archFileName  = ".go-arch-lint.yml"
	goModFileName = "go.mod"
)

func Execute() {
	setDefaultCommandIfNonePresent("--help")

	defer func() {
		if err := recover(); err != nil {
			if goErr, ok := err.(error); ok {
				halt(fmt.Errorf("panic: %s", goErr))
			} else {
				halt(fmt.Errorf("panic: %s", err))
			}
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		halt(fmt.Errorf("error: %s", err))
	}

	os.Exit(0)
}

func setDefaultCommandIfNonePresent(defaultCommand string) {
	if len(os.Args) != 1 {
		return
	}

	os.Args = append(os.Args, defaultCommand)
}

func init() {
	// color-output
	rootCmd.PersistentFlags().BoolVar(
		&flagUseColors,
		flagNameOutputColors,
		true,
		"use ANSI colors in terminal output",
	)

	// color-output
	rootCmd.PersistentFlags().BoolVar(
		&flagOutputJsonOneLine,
		flagNameOutputJsonOneLine,
		false,
		"format JSON as single line payload (without line breaks), only for json output type",
	)

	// output-type
	rootCmd.PersistentFlags().StringVar(
		&flagOutputType,
		flagNameOutputType,
		outputTypeDefault,
		fmt.Sprintf("type of command output, variants: [%s]", strings.Join(outputTypeVariants, ", ")),
	)
	rootCmd.PersistentFlags().BoolVar(
		&flagAliasJson,
		flagNameAliasJson,
		false,
		fmt.Sprintf("(alias for --%s=%s)",
			flagNameOutputType,
			outputTypeJSON,
		),
	)

	// max warnings
	rootCmd.PersistentFlags().IntVar(
		&flagMaxWarnings,
		flagNameMaxWarnings,
		512,
		"max number of warnings to output",
	)

	// project-path
	rootCmd.PersistentFlags().StringVar(
		&flagProjectPath,
		flagNameProjectPath,
		"",
		fmt.Sprintf("absolute path to project directory (where '%s' is located)", archFileName),
	)

	// init
	cobra.OnInitialize(func() {
		au = aurora.NewAurora(flagUseColors)

		processAliases()
		assertFlagOutputTypeValid()
		assertFlagMaxWarningsValid()
	})
}

func processAliases() {
	processAliasJson()
}

func processAliasJson() {
	if !flagAliasJson {
		return
	}

	if flagOutputType != outputTypeDefault && flagOutputType != outputTypeJSON {
		panic(fmt.Sprintf("flag --%s not compatible with --%s",
			flagNameAliasJson,
			flagNameOutputType,
		))
	}

	flagOutputType = outputTypeJSON
}

func assertFlagOutputTypeValid() {
	if flagOutputType == outputTypeDefault {
		flagOutputType = outputTypeASCII
	}

	for _, variant := range outputTypeVariants {
		if flagOutputType == variant {
			return
		}
	}

	panic(fmt.Sprintf("unknown output-type: %s", flagOutputType))
}

func assertFlagMaxWarningsValid() {
	const warningsRangeMin = 1
	const warningsRangeMax = 32768

	if flagMaxWarnings < warningsRangeMin || flagMaxWarnings > warningsRangeMax {
		panic(fmt.Sprintf("flag '%s' should by in range [%d .. %d]",
			flagNameMaxWarnings,
			warningsRangeMin,
			warningsRangeMax,
		))
	}
}
