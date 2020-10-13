package validator

import (
	"fmt"
)

func withCheckerVendorDependencies(reg checkerRegistry) {
	for name, rules := range reg.spec().Dependencies {
		name := name
		rules := rules

		existVendors := make(map[string]bool)

		for index, vendorName := range rules.CanUse {
			vendorName := vendorName

			if _, ok := existVendors[vendorName]; ok {
				reg.applyChecker(
					fmt.Sprintf("$.deps.%s.canUse[%d]", name, index),
					func() error {
						return fmt.Errorf("vendor dublicated")
					},
				)
			}

			reg.applyChecker(
				fmt.Sprintf("$.deps.%s.canUse[%d]", name, index),
				func() error {
					return reg.utils().isKnownVendor(vendorName)
				},
			)

			existVendors[vendorName] = true
		}
	}
}
