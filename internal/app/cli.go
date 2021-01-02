package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/internal/app/internal/container"
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/version"
)

func Execute() int {
	di := container.NewContainer(
		version.Version,
		version.BuildTime,
		version.CommitHash,
	)
	rootCmd := di.ProvideRootCommand()

	err := rootCmd.Execute()
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
