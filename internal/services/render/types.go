package render

import "github.com/fe3dback/go-arch-lint/internal/models"

type (
	referenceRender interface {
		SourceCode(ref models.CodeReference, highlight bool) []byte
	}

	colorPrinter interface {
		Red(in string) (out string)
		Green(in string) (out string)
		Yellow(in string) (out string)
		Blue(in string) (out string)
		Magenta(in string) (out string)
		Cyan(in string) (out string)
		Gray(in string) (out string)
	}
)
