package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type ComponentsValidator struct {
	pathHelper pathHelper
}

func NewComponentsValidator(
	pathHelper pathHelper,
) *ComponentsValidator {
	return &ComponentsValidator{
		pathHelper: pathHelper,
	}
}

func (c *ComponentsValidator) Validate(ctx *validationContext) {
	if ctx.conf.Components.Map.Len() == 0 {
		ctx.AddNotice(
			"at least one component should by defined",
			ctx.conf.Version.Ref,
		)

		return
	}

	ctx.conf.Components.Map.Each(func(_ models.ComponentName, component models.ConfigComponent, _ models.Reference) {
		for _, pathGlob := range component.In {
			matched, err := c.pathHelper.FindProjectFiles(models.FileQuery{
				Path:             pathGlob.Value,
				WorkingDirectory: ctx.conf.WorkingDirectory.Value,
				Type:             models.FileMatchQueryTypeOnlyDirectories,
			})
			if err != nil {
				ctx.AddNotice(
					fmt.Sprintf("failed find directories: %v", err),
					pathGlob.Ref,
				)
				return
			}

			if len(matched) == 0 {
				ctx.AddNotice(
					fmt.Sprintf("not found any directories by glob '%s'", pathGlob.Value),
					pathGlob.Ref,
				)
				return
			}
		}
	})
}
