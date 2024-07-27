package validator

import (
	"fmt"
)

type WorkdirValidator struct {
	pathHelper pathHelper
}

func NewWorkdirValidator(
	pathHelper pathHelper,
) *WorkdirValidator {
	return &WorkdirValidator{
		pathHelper: pathHelper,
	}
}

func (c *WorkdirValidator) Validate(ctx *validationContext) {
	workDir := ctx.conf.WorkingDirectory.Value
	matched, err := c.pathHelper.MatchProjectFiles(workDir, true)
	if err != nil {
		ctx.AddNotice(
			fmt.Sprintf("failed find directory '%s': %v", workDir, err),
			ctx.conf.WorkingDirectory.Ref,
		)
		return
	}

	if len(matched) == 0 {
		ctx.SkipOtherValidators()
		ctx.AddNotice(
			fmt.Sprintf("not found directory '%s', possible not exist", workDir),
			ctx.conf.WorkingDirectory.Ref,
		)
		return
	}
}
