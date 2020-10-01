package validator

import "fmt"

func withCheckerVendors(reg checkerRegistry) {
	for name, vendor := range reg.spec().Vendors {
		vendor := vendor

		reg.applyChecker(
			fmt.Sprintf("$.vendors.%s.in", name),
			func() error {
				return reg.utils().isValidImportPath(vendor.ImportPath)
			},
		)
	}
}
