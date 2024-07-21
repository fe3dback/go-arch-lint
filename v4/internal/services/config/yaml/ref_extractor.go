package yaml

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

// this will parse custom yaml error with source annotation
// and extract File/Line/Pos info to common Reference struct
//
// example:
//
//	||       6 |   imports:
//	||       7 |     strictMode: false
//	||       8 |     allowAnyVendorImports: true
//	||    >  9 |   structTags:
//	||             ^
//	||      10 |     allowed: true
//	||      11 |     # allowed: true                # all tags (default)
//	||      12 |     # allowed: false               # no tags
//	||      13 |
func extractReferenceFromError(tCtx TransformContext, err error) models.Reference {
	// todo:
	// yaml.FormatError()
	return models.NewReference(tCtx.file, 1, 0)
}
