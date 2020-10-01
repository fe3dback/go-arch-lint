package annotated_validator

import "github.com/fe3dback/go-arch-lint/models"

type (
	AnnotatedWarningParser interface {
		Parse(sourceText string) (line, pos int, err error)
	}

	Validator interface {
		Validate() []models.ArchFileSyntaxWarning
	}
)
