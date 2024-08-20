package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/v4/internal/app/internal/container"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func Execute() (exitCode int) {
	di := container.NewContainer()
	cliApp := di.Cli()

	err := cliApp.Run(os.Args)
	if err != nil {
		userError := &models.UserLandError{}
		if !errors.As(err, &userError) {
			fmt.Printf("Error: %s\n", err)
		}

		return 1
	}

	return 0
}
