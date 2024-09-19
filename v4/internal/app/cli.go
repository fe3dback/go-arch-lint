package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/v4/internal/app/internal/container"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

const (
	EnvTypeMain = "main"

	// EnvTypeTests is special env for tests. Main different is that
	// Execute function will be called multiple times, with different os.args
	EnvTypeTests = "tests"
)

type (
	EnvType string
)

func Execute(env EnvType) (exitCode int) {
	di := container.NewContainer()
	if env == EnvTypeTests {
		// unit tests will call Execute many times with different input args
		// so, we need build new DI heap-cache every run
		di.DropBuiltHeap()
	}

	cliApp := di.Cli()

	err := cliApp.Run(os.Args)
	if err != nil {
		userLandError := &models.UserLandError{}

		switch {
		case errors.Is(err, models.ProjectNotMatchSpecificationsError):
			// actually is not a linter error, just enforce exitCode=1
		case errors.As(err, &userLandError):
			// this error will be printed with custom formatted view
		default:
			// internal errors we can just print as-is
			fmt.Printf("Error: %s\n", err)
		}

		return 1
	}

	return 0
}
