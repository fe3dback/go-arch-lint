package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmdtest"
)

const binaryName = "go-arch-lint"

var update = flag.Bool("update", false, "update test files with results")

func TestCLI(t *testing.T) {
	ts, err := cmdtest.Read("tests/**")
	if err != nil {
		t.Fatal(err)
	}

	ts.Setup = func(_ string) error {
		_, testFileName, _, ok := runtime.Caller(0)
		if !ok {
			return fmt.Errorf("failed get real working directory from caller")
		}

		projectRootDir := filepath.Dir(testFileName)
		if err := os.Setenv("ROOTDIR", projectRootDir); err != nil {
			return fmt.Errorf("failed change 'ROOTDIR' to caller working directory: %w", err)
		}

		// additional flag for github.com/urfave/cli/v2
		// will print errors in some cases (instead of silent fail)
		_ = os.Setenv("CLI_TEMPLATE_ERROR_DEBUG", "1")

		return nil
	}

	ts.Commands[binaryName] = cmdtest.InProcessProgram(binaryName, runWithinTests)
	ts.Run(t, *update)
}
