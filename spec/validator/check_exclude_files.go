package validator

import (
	"fmt"
	"regexp"
)

func withCheckerExcludedFiles(reg checkerRegistry) {
	for index, regExp := range reg.spec().ExcludeFilesRegExp {
		reg.applyChecker(
			fmt.Sprintf("$.excludeFiles[%d]", index),
			func() error {
				_, err := regexp.Compile(regExp)
				return err
			},
		)
	}
}
