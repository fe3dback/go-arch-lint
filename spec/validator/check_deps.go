package validator

import (
	"fmt"
)

func withCheckerDependencies(reg checkerRegistry) {
	for name, rules := range reg.spec().Dependencies {
		name := name
		rules := rules

		reg.applyChecker(fmt.Sprintf("$.deps.%s", name), func() error {
			return reg.utils().isKnownComponent(name)
		})

		if len(rules.MayDependOn) == 0 && len(rules.CanUse) == 0 {
			reg.applyChecker(
				fmt.Sprintf("$.deps.%s", name),
				func() error {
					if rules.AnyProjectDeps {
						return nil
					}

					if rules.AnyVendorDeps {
						return nil
					}

					return fmt.Errorf("should have ref in 'MayDependOn' or at least one flag of ['anyProjectDeps', 'anyVendorDeps']")
				},
			)
		}
	}
}
