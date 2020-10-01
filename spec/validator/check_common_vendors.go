package validator

import "fmt"

func withCheckerCommonVendors(reg checkerRegistry) {
	for index, vendorName := range reg.spec().CommonVendors {
		vendorName := vendorName

		reg.applyChecker(
			fmt.Sprintf("$.commonVendors[%d]", index),
			func() error {
				return reg.utils().isKnownVendor(vendorName)
			},
		)
	}
}
