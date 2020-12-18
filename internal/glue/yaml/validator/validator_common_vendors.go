package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorCommonVendors struct {
	refResolver YamlSourceCodeReferenceResolver
	utils       *utils
}

func newValidatorCommonVendors(
	refResolver YamlSourceCodeReferenceResolver,
	utils *utils,
) *validatorCommonVendors {
	return &validatorCommonVendors{
		refResolver: refResolver,
		utils:       utils,
	}
}

func (v *validatorCommonVendors) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for index, vendorName := range doc.CommonVendors {
		if err := v.utils.assertKnownVendor(vendorName); err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    v.refResolver.Resolve(fmt.Sprintf("$.commonVendors[%d]", index)),
			})
		}
	}

	return notices
}
