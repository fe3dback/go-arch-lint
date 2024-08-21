package container

import (
	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func (c *Container) sdk() *sdk.SDK {
	return once(func() *sdk.SDK {
		return sdk.NewSDK(
			arch.PathAbsolute(c.cCtx.String(models.FlagProjectPath)),
			sdk.WithUsedContext(arch.UsedContextCLI),
			sdk.WithSkipMissUse(c.cCtx.Bool(models.FlagSkipMissUsages)),
		)
	})
}
