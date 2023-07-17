package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
)

type validatorDeps struct {
	utils *utils
}

func newValidatorDeps(
	utils *utils,
) *validatorDeps {
	return &validatorDeps{
		utils: utils,
	}
}

func (v *validatorDeps) Validate(doc spec.Document) []arch.Notice {
	notices := make([]arch.Notice, 0)

	for name, rule := range doc.Dependencies() {
		if err := v.utils.assertKnownComponent(name); err != nil {
			notices = append(notices, arch.Notice{
				Notice: err,
				Ref:    rule.Reference,
			})
		}

		if len(rule.Value.MayDependOn()) > 0 && rule.Value.AnyProjectDeps().Value {
			notices = append(notices, arch.Notice{
				Notice: fmt.Errorf("'anyProjectDeps=true' used with not empty 'MayDependOn' list (likely this is miss configuration)"),
				Ref:    rule.Value.AnyProjectDeps().Reference,
			})
		}

		if len(rule.Value.CanUse()) > 0 && rule.Value.AnyVendorDeps().Value {
			notices = append(notices, arch.Notice{
				Notice: fmt.Errorf("'AnyVendorDeps=true' used with not empty 'CanUse' list (likely this is miss configuration)"),
				Ref:    rule.Value.AnyVendorDeps().Reference,
			})
		}

		if len(rule.Value.MayDependOn()) == 0 && len(rule.Value.CanUse()) == 0 {
			if rule.Value.AnyProjectDeps().Value {
				continue
			}

			if rule.Value.AnyVendorDeps().Value {
				continue
			}

			notices = append(notices, arch.Notice{
				Notice: fmt.Errorf("should have ref in 'mayDependOn'/'canUse' or at least one flag of ['anyProjectDeps', 'anyVendorDeps']"),
				Ref:    rule.Reference,
			})
		}
	}

	return notices
}
