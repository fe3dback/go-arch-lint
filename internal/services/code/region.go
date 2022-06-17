package code

import "github.com/fe3dback/go-arch-lint/internal/models"

type (
	codeRegion struct {
		lineFirst int
		lineMain  int
		lineLast  int
	}
)

func calculateCodeRegion(ref models.CodeReference, maxLines int) codeRegion {
	lineFirst := ref.LineFrom
	lineMain := ref.Pointer.Line
	lineLast := ref.LineTo

	if lineFirst > lineLast {
		lineFirst, lineLast = lineLast, lineFirst
	}

	if lineFirst < 1 {
		lineFirst = 1
	}

	if lineLast > maxLines {
		lineLast = maxLines
	}

	if lineMain < lineFirst {
		lineMain = lineFirst
	}

	if lineMain > lineLast {
		lineMain = lineLast
	}

	return codeRegion{
		lineFirst: lineFirst,
		lineMain:  lineMain,
		lineLast:  lineLast,
	}
}
