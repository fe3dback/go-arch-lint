package container

import (
	"fmt"

	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func (c *Container) sdk() *sdk.SDK {
	return once(func() *sdk.SDK {
		createdSDK, err := sdk.NewSDK(
			arch.PathAbsolute(c.cCtx.String(models.FlagProjectPath)),
			sdk.WithUsedContext(arch.UsedContextCLI),
			sdk.WithSkipMissUse(c.cCtx.Bool(models.FlagSkipMissUsages)),
			sdk.WithOutputColors(c.cCtx.Bool(models.FlagOutputUseAsciiColors)),
		)
		if err != nil {
			panic(fmt.Errorf("failed to initialize sdk: %w", err))
		}

		return createdSDK
	})
}
