package main

import (
	"os"

	"github.com/fe3dback/go-arch-lint/v4/internal/app"
)

func main() {
	os.Exit(run())
}

func run() int {
	return app.Execute(app.EnvTypeMain)
}

func runWithinTests() int {
	return app.Execute(app.EnvTypeTests)
}
