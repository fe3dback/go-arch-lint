package check

import (
	"fmt"

	"github.com/urfave/cli/v2"

	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/commands/check"
)

type Command struct {
	sdk         *sdk.SDK
	specFetcher specFetcher
}

func NewCommand(sdk *sdk.SDK, specFetcher specFetcher) *Command {
	return &Command{
		sdk:         sdk,
		specFetcher: specFetcher,
	}
}

func (c *Command) Execute(cCtx *cli.Context) (any, error) {
	in := c.parseIn(cCtx)
	err := c.validateIn(in)
	if err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	spec, err := c.specFetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spec: %w", err)
	}

	out, err := c.sdk.Check(spec, in)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (c *Command) parseIn(cCtx *cli.Context) check.In {
	in := check.In{}
	in.MaxWarnings = cCtx.Int(flagMaxWarnings)

	return in
}

func (c *Command) validateIn(in check.In) error {
	const warningsRangeMin = 1
	const warningsRangeMax = 32768

	if in.MaxWarnings < warningsRangeMin || in.MaxWarnings > warningsRangeMax {
		return fmt.Errorf(
			"flag '%s' should by in range [%d .. %d]",
			flagMaxWarnings,
			warningsRangeMin,
			warningsRangeMax,
		)
	}

	return nil
}
