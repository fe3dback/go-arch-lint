package validator

import (
	"fmt"
	"regexp"
)

type ExcludedFilesValidator struct{}

func NewExcludedFilesValidator() *ExcludedFilesValidator {
	return &ExcludedFilesValidator{}
}

func (c *ExcludedFilesValidator) Validate(ctx *validationContext) {
	for _, excludedFile := range ctx.conf.Exclude.RelativeFiles {
		_, err := regexp.Compile(string(excludedFile.Value))
		if err != nil {
			ctx.SkipOtherValidators()
			ctx.AddNotice(
				fmt.Sprintf("Regexp pattern '%s' for files exclusion is invalid", excludedFile.Value),
				excludedFile.Ref,
			)
		}
	}
}
