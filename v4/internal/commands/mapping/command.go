package mapping

import (
	"fmt"

	"github.com/urfave/cli/v2"

	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/mapping"
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

	spec, err := c.specFetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spec: %w", err)
	}

	out, err := c.sdk.Mapping(spec, in)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (c *Command) parseIn(cCtx *cli.Context) mapping.In {
	in := mapping.In{}
	in.Scheme = cCtx.String(flagScheme)

	return in
}
