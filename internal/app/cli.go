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
		version.VERSION,
		version.BUILD_TIME,
		version.COMMIT_HASH,
	)
	rootCmd := di.ProvideRootCommand()

	err := rootCmd.Execute()
	if err != nil {
		if errors.Is(err, models.UserSpaceError{}) {
			// do not display user space errors (usually explain will by in ascii/json output)
			return 1
		}

		// system error, not possible to output this in json, so just dump to stdout
		_, _ = fmt.Fprintln(os.Stderr, fmt.Errorf("%s\n", err))
		return 1
	}

	return 0
}
