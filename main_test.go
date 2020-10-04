package main

import (
	"flag"
	"testing"

	"github.com/google/go-cmdtest"
)

const binaryName = "go-arch-lint"

var update = flag.Bool("update", false, "update test files with results")

func TestCLI(t *testing.T) {
	ts, err := cmdtest.Read("test/**")
	if err != nil {
		t.Fatal(err)
	}

	ts.Commands[binaryName] = cmdtest.InProcessProgram(binaryName, run)
	ts.Run(t, *update)
}
