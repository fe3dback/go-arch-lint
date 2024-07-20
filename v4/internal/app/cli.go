package app

import (
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/v4/internal/app/internal/container"
)

func Execute() int {
	di := container.NewContainer()
	cliApp := di.Cli()

	err := cliApp.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %s\n", err)

		return 1
	}

	return 0
}
