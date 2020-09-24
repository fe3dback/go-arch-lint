package validator

import "fmt"

func withCheckerCommonComponents(reg checkerRegistry) {
	for index, componentName := range reg.spec().CommonComponents {
		reg.applyChecker(
			fmt.Sprintf("$.commonComponents[%d]", index),
			func() error {
				return reg.utils().isKnownComponent(componentName)
			},
		)
	}
}
