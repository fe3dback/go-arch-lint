package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type validatorVendors struct {
	refResolver YamlSourceCodeReferenceResolver
	utils       *utils
}

func newValidatorVendors(
	refResolver YamlSourceCodeReferenceResolver,
	utils *utils,
) *validatorVendors {
	return &validatorVendors{
		refResolver: refResolver,
		utils:       utils,
	}
}

func (v *validatorVendors) Validate(doc spec.Document) []speca.Notice {
	notices := make([]speca.Notice, 0)

	for name, vendor := range doc.Vendors {
		err := v.utils.assertVendorImportPathValid(vendor.ImportPath)
		if err != nil {
			notices = append(notices, speca.Notice{
				Notice: err,
				Ref:    v.refResolver.Resolve(fmt.Sprintf("$.vendors.%s.in", name)),
			})
		}
	}

	return notices
}
