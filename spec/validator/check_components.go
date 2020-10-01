package validator

import (
	"fmt"
)

func withCheckerComponents(reg checkerRegistry) {
	reg.applyChecker("$.components", func() error {
		if len(reg.spec().Components) == 0 {
			return fmt.Errorf("at least one component should by defined")
		}

		return nil
	})

	for name, component := range reg.spec().Components {
		component := component

		reg.applyChecker(
			fmt.Sprintf("$.components.%s.in", name),
			func() error {
				return reg.utils().isValidPath(component.LocalPath)
			},
		)
	}
}
