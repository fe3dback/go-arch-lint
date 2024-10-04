package xstdout

import (
	"errors"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/codeprinter"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type ErrorBuilder struct {
	codePrinter codePrinter
}

func NewErrorBuilder(codePrinter codePrinter) *ErrorBuilder {
	return &ErrorBuilder{
		codePrinter: codePrinter,
	}
}

func (eb *ErrorBuilder) BuildError(err error) models.CmdStdoutErrorOut {
	noticeErr := &arch.ErrorWithNotices{}
	if errors.As(err, &noticeErr) {
		return models.CmdStdoutErrorOut{
			OverallNote: noticeErr.OverallMessage,
			Errors:      noticeErr.Notices,
		}
	}

	refError := arch.ReferencedError{}
	if errors.As(err, &refError) {
		return eb.buildSingleRefError(refError)
	}

	return models.CmdStdoutErrorOut{
		Errors: []arch.Notice{
			{
				Message:   err.Error(),
				Reference: arch.NewInvalidReference(),
			},
		},
	}
}

func (eb *ErrorBuilder) buildSingleRefError(refError arch.ReferencedError) models.CmdStdoutErrorOut {
	ref := refError.Reference()

	preview := ""
	if ref.Valid {
		preview, _ = eb.codePrinter.Print(transformRef(ref), codeprinter.CodePrintOpts{
			LineNumbers: true,
			Arrows:      true,
			Mode:        codeprinter.CodePrintModeExtend,
		})
	}

	return models.CmdStdoutErrorOut{
		Errors: []arch.Notice{
			{
				Message:     refError.Error(),
				Reference:   ref,
				CodePreview: preview,
			},
		},
	}
}

func transformRef(ref arch.Reference) codeprinter.Reference {
	return codeprinter.Reference{
		File:   string(ref.File),
		Line:   ref.Line,
		Column: ref.Column,
		Valid:  ref.Valid,
	}
}
