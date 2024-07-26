package validator

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type validationContext struct {
	conf    models.Config
	notices []models.Notice
}

func (vc *validationContext) AddNotice(message string, ref models.Reference) {
	vc.notices = append(vc.notices, models.Notice{
		Message:   message,
		Reference: ref,
	})
}

func (vc *validationContext) IsKnownComponent(name models.ComponentName) bool {
	return vc.conf.Components.Map.Has(name)
}

func (vc *validationContext) IsKnownVendor(name models.VendorName) bool {
	return vc.conf.Vendors.Map.Has(name)
}
