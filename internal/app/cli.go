package app

import (
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint/internal/app/internal/container"
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
		_, _ = fmt.Fprintln(os.Stderr, fmt.Errorf("failed: %s", err))

		return 1
	}

	return 0
}
