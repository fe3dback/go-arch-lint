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
)

var (
	// flags
	flagUseColors         bool
	flagOutputType        outputType
	flagOutputJsonOneLine bool
	flagMaxWarnings       int
	flagProjectPath       string

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
	outputTypeASCII = "ascii"
	outputTypeJSON  = "json"
)

const (
	archFileName  = ".go-arch-lint.yml"
	goModFileName = "go.mod"
)

func Execute() {
	setDefaultCommandIfNonePresent(rootCmd, cmdNameCheck)

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

func setDefaultCommandIfNonePresent(rootCmd *cobra.Command, defaultCommand string) {
	if len(os.Args) == 0 {
		return
	}

	for _, arg := range os.Args {
		for _, command := range rootCmd.Commands() {
			if command.Use == arg {
				return
			}
		}
	}

	os.Args = append([]string{os.Args[0], defaultCommand}, os.Args[1:]...)
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
		"format JSON as multiline payload, only for json output type",
	)

	// output-type
	rootCmd.PersistentFlags().StringVar(
		&flagOutputType,
		flagNameOutputType,
		outputTypeASCII,
		fmt.Sprintf("type of command output, variants: [%s]", strings.Join(outputTypeVariants, ", ")),
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

		assertFlagOutputTypeValid()
		assertFlagMaxWarningsValid()
	})
}

func assertFlagOutputTypeValid() {
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
