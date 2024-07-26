package validator

import (
	"fmt"
	"path"

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

	ctx.conf.Components.Map.Each(func(name models.ComponentName, component models.ConfigComponent, reference models.Reference) {
		for _, pathGlob := range component.In {
			relPath := models.PathRelativeGlob(path.Join(string(ctx.conf.WorkingDirectory.Value), string(pathGlob.Value)))
			matched, err := c.pathHelper.MatchProjectFiles(relPath)
			if err != nil {
				ctx.AddNotice(
					fmt.Sprintf("failed find files: %v", err),
					pathGlob.Ref,
				)
				return
			}

			if len(matched) == 0 {
				ctx.AddNotice(
					fmt.Sprintf("not found any files by glob '%s'", pathGlob.Value),
					pathGlob.Ref,
				)
				return
			}
		}
	})
}
