package validator

import (
	"fmt"
)

const supportedVersion = 1

func withCheckerVersion(reg checkerRegistry) {
	reg.applyChecker("$.version", func() error {
		if reg.spec().Version <= supportedVersion && reg.spec().Version > 0 {
			return nil
		}

		return fmt.Errorf("version %d is not supported, supported: [%d]",
			reg.spec().Version,
			supportedVersion,
		)
	})
}
