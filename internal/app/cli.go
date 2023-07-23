package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/internal/app/internal/container"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

func Execute() int {
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// -- build DI
	di := container.NewContainer(
		Version,
		BuildTime,
		CommitHash,
	)

	// -- process
	err := di.CommandRoot().ExecuteContext(mainCtx)

	// -- handle errors
	if err != nil {
		if errors.Is(err, models.UserSpaceError{}) {
			// do not display user space errors (usually explain will be in ascii/json output)
			return 1
		}

		// system error, not possible to output this in json, so just dump to stdout
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return 1
	}

	return 0
}
