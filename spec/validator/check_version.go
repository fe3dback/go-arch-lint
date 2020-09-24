package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/spec/archfile"
)

func withCheckerVersion(reg checkerRegistry) {
	reg.applyChecker("$.version", func() error {
		if reg.spec().Version <= archfile.SupportedVersion && reg.spec().Version > 0 {
			return nil
		}

		return fmt.Errorf("version %d is not supported, supported: [%d]",
			reg.spec().Version,
			archfile.SupportedVersion,
		)
	})
}
