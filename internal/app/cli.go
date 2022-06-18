package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	terminal "github.com/fe3dback/span-terminal"

	"github.com/fe3dback/go-arch-lint/internal/app/internal/container"
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/version"
)

func Execute() int {
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	terminal.SetGlobalTerminal(terminal.NewTerminal(
		terminal.WithStdoutMaxLines(6),
		terminal.WithContainerMaxLines(3),
		terminal.WithRenderOpts(
			terminal.WithRenderOptSpanMaxRoots(4),
			terminal.WithRenderOptSpanMaxChild(8),
			terminal.WithRenderOptSpanMaxDetails(16),
		),
	))

	// -- build DI
	di := container.NewContainer(
		version.Version,
		version.BuildTime,
		version.CommitHash,
	)

	// -- process
	rootCmd := di.ProvideRootCommand()
	err := rootCmd.ExecuteContext(mainCtx)

	// -- handle errors
	if err != nil {
		if errors.Is(err, models.UserSpaceError{}) {
			// do not display user space errors (usually explain will by in ascii/json output)
			return 1
		}

		// system error, not possible to output this in json, so just dump to stdout
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return 1
	}

	return 0
}
