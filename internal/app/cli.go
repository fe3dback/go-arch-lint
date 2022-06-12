package app

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	terminal "github.com/fe3dback/span-terminal"

	"github.com/fe3dback/go-arch-lint/internal/app/internal/container"
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/version"
)

func Execute() int {
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	frontendTerminal := terminal.NewTerminal(mainCtx, os.Stdout)

	defaultStdout := os.Stdout
	logBufferReader, logBufferWriter, _ := os.Pipe()

	// forward all logs and output to buffer
	os.Stdout = logBufferWriter
	log.SetOutput(logBufferWriter)

	// register frontend terminal
	terminal.RegisterTerminal(frontendTerminal)

	// -- handle signals

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Kill, os.Interrupt)

		_ = <-sig
		cancel()
	}()

	// -- build DI
	di := container.NewContainer(
		version.Version,
		version.BuildTime,
		version.CommitHash,
	)

	// -- process
	rootCmd := di.ProvideRootCommand()
	err := rootCmd.ExecuteContext(mainCtx)

	// -- write all logs
	terminal.Shutdown()
	time.Sleep(time.Millisecond * 500)

	// -- write buffered logs
	_ = logBufferWriter.Close()
	bufferedOutput, _ := ioutil.ReadAll(logBufferReader)
	os.Stdout = defaultStdout

	fmt.Printf("%s", bufferedOutput)

	// -- handle errors
	if err != nil {
		if errors.Is(err, models.UserSpaceError{}) {
			// do not display user space errors (usually explain will by in ascii/json output)
			return 1
		}

		// system error, not possible to output this in json, so just dump to stdout
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return 1
	}

	return 0
}
