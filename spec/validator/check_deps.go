package validator

import (
	"fmt"
)

func withCheckerDependencies(reg checkerRegistry) {
	for name, rules := range reg.spec().Dependencies {
		reg.applyChecker(fmt.Sprintf("$.deps.%s", name), func() error {
			return reg.utils().isKnownComponent(name)
		})

		for index, componentName := range rules.MayDependOn {
			reg.applyChecker(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, index), func() error {
				return reg.utils().isKnownComponent(componentName)
			})
		}

		for index, vendorName := range rules.CanUse {
			reg.applyChecker(
				fmt.Sprintf("$.deps.%s.canUse[%d]", name, index),
				func() error {
					return reg.utils().isKnownVendor(vendorName)
				},
			)
		}

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
