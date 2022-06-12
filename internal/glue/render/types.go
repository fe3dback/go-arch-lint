package render

import "github.com/fe3dback/go-arch-lint/internal/models"

type (
	ReferenceRender interface {
		SourceCode(ref models.Reference, height int, highlight bool) []byte
	}

	ColorPrinter interface {
		Red(in string) (out string)
		Green(in string) (out string)
		Yellow(in string) (out string)
		Blue(in string) (out string)
		Magenta(in string) (out string)
		Cyan(in string) (out string)
		Gray(in string) (out string)
	}
)
