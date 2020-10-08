package main

import (
	"os"

	"github.com/fe3dback/go-arch-lint/cmd"
)

func main() {
	os.Exit(run())
}

func run() int {
	return cmd.Execute()
}
