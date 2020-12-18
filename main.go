package main

import (
	"os"

	"github.com/fe3dback/go-arch-lint/internal/app"
)

func main() {
	os.Exit(run())
}

func run() int {
	return app.Execute()
}
