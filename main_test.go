package main

import (
	"testing"

	"github.com/google/go-cmdtest"
)

func TestCLI(t *testing.T) {
	ts, err := cmdtest.Read("test")
	if err != nil {
		t.Fatal(err)
	}

	ts.Commands["go-arch-lint"] = cmdtest.InProcessProgram("my-cli", run)
	ts.Run(t, false)
}
