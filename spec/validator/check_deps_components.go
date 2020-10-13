package validator

import (
	"fmt"
)

func withCheckerComponentDependencies(reg checkerRegistry) {
	for name, rules := range reg.spec().Dependencies {
		name := name
		rules := rules

		existComponents := make(map[string]bool)

		for index, componentName := range rules.MayDependOn {
			componentName := componentName

			if _, ok := existComponents[componentName]; ok {
				reg.applyChecker(
					fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, index),
					func() error {
						return fmt.Errorf("component dublicated")
					},
				)
			}

			reg.applyChecker(
				fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, index),
				func() error {
					return reg.utils().isKnownComponent(componentName)
				},
			)

			existComponents[componentName] = true
		}
	}
}
