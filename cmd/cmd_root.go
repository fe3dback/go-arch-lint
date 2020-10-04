package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/fe3dback/go-arch-lint/models"
)

// command names
const (
	cmdNameCheck   = "check"
	cmdNameVersion = "version"
)

// ctx storage key names
const (
	ctxRootInputFlags = "rootFlags"
)

// global defaults
const (
	goModFileName = "go.mod"
)

var (
	outputTypeVariantsConst = []string{outputTypeASCII, outputTypeJSON}
)

type (
	outputType = string
)

func Execute() (exitCode int) {
	// prepare flags
	defaults := newDefaultFlags()

	// prepare os args
	setDefaultCommandIfNonePresent("help")

	// catch next errors
	defer func() {
		if err := recover(); err != nil {
			// turn off colored output for compatible
			defaults.useColors = false
			exitCode = 1

			// local panic (without trace)
			if goErr, ok := err.(models.UserSpaceError); ok {
				halt(defaults, fmt.Errorf("panic: %s", goErr))
				return
			}

			// output panic + trace
			halt(defaults, fmt.Errorf(
				"panic: %s\n--\n%s",
				err,
				debug.Stack(),
			))
			return
		}
	}()

	// assemble commands
	rootCmd := assembleRootCommand()

	// parse flags
	flags := parseFlags(rootCmd, defaults)
	ctx := context.WithValue(context.Background(), ctxRootInputFlags, flags)

	// execute
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		halt(flags, fmt.Errorf("error: %s", err))
		return 1
	}

	// ok
	return 0
}

func setDefaultCommandIfNonePresent(defaultCommand string) {
	if len(os.Args) != 1 {
		return
	}

	os.Args = append(os.Args, defaultCommand)
}
