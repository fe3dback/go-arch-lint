package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
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

func (v *validatorDeps) Validate(doc arch.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for name, rules := range doc.Dependencies().Map() {
		if err := v.utils.assertKnownComponent(name); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    rules.Reference(),
			})
		}

		if len(rules.MayDependOn()) == 0 && len(rules.CanUse()) == 0 {
			if rules.AnyProjectDeps().Value() {
				continue
			}

			if rules.AnyVendorDeps().Value() {
				continue
			}

			notices = append(notices, speca.Notice{
				Notice: fmt.Errorf("should have ref in 'mayDependOn' or at least one flag of ['anyProjectDeps', 'anyVendorDeps']"),
				Ref:    rules.Reference(),
			})
		}
	}

	return notices
}
