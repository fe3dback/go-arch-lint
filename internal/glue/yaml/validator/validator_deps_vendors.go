package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorDepsVendors struct {
	refResolver YamlSourceCodeReferenceResolver
	utils       *utils
}

func newValidatorDepsVendors(
	refResolver YamlSourceCodeReferenceResolver,
	utils *utils,
) *validatorDepsVendors {
	return &validatorDepsVendors{
		refResolver: refResolver,
		utils:       utils,
	}
}

func (v *validatorDepsVendors) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for name, rules := range doc.Dependencies {
		existVendors := make(map[string]bool)

		for index, vendorName := range rules.CanUse {
			refPath := fmt.Sprintf("$.deps.%s.canUse[%d]", name, index)

			if _, ok := existVendors[vendorName]; ok {
				notices = append(notices, speca.Notice{
					Notice: fmt.Errorf("vendor '%s' dublicated in '%s' deps", vendorName, name),
					Ref:    v.refResolver.Resolve(refPath),
				})
			}

			if err := v.utils.assertKnownVendor(vendorName); err != nil {
				notices = append(notices, speca.Notice{
					Notice: err,
					Ref:    v.refResolver.Resolve(refPath),
				})
			}

			existVendors[vendorName] = true
		}
	}

	return notices
}
