package container

import (
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/project/xpath"
)

func (c *Container) serviceProjectPathHelper() *xpath.Helper {
	return once(func() *xpath.Helper {
		return xpath.NewHelper(
			c.cCtx.String(models.FlagProjectPath),
		)
	})
}
