package spec

import (
	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type Fetcher struct {
	sdk        *sdk.SDK
	configPath arch.PathRelative
}

func NewFetcher(
	sdk *sdk.SDK,
	configPath arch.PathRelative,
) *Fetcher {
	return &Fetcher{
		sdk:        sdk,
		configPath: configPath,
	}
}

func (f *Fetcher) Fetch() (arch.Spec, error) {
	return f.sdk.Spec().FromFile(f.configPath)
}
