package xstdout

import (
	"errors"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/codeprinter"
	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type ErrorBuilder struct {
	codePrinter codePrinter
	useColors   bool
}

func NewErrorBuilder(codePrinter codePrinter, useColors bool) *ErrorBuilder {
	return &ErrorBuilder{
		codePrinter: codePrinter,
		useColors:   useColors,
	}
}

func (eb *ErrorBuilder) BuildError(err error) models.CmdStdoutErrorOut {
	noticeErr := &arch.ErrorWithNotices{}
	if errors.As(err, &noticeErr) {
		return eb.buildNoticesError(noticeErr)
	}

	refError := arch.ReferencedError{}
	if errors.As(err, &refError) {
		return eb.buildSingleError(refError)
	}

	return models.CmdStdoutErrorOut{
		Errors: []models.StdoutNotice{
			{
				Notice: arch.Notice{
					Message:   err.Error(),
					Reference: arch.NewInvalidReference(),
				},
			},
		},
	}
}

func (eb *ErrorBuilder) buildNoticesError(err *arch.ErrorWithNotices) models.CmdStdoutErrorOut {
	outNotices := make([]models.StdoutNotice, 0, len(err.Notices))

	mode := codeprinter.CodePrintModeExtend
	if len(err.Notices) > 10 {
		mode = codeprinter.CodePrintModeOneLine
	}

	for _, notice := range err.Notices {
		preview := ""
		if notice.Reference.Valid {
			preview, _ = eb.codePrinter.Print(transformRef(notice.Reference), codeprinter.CodePrintOpts{
				LineNumbers: true,
				Arrows:      true,
				Highlight:   eb.useColors,
				Mode:        mode,
			})
		}

		outNotices = append(outNotices, models.StdoutNotice{
			Notice:  notice,
			Preview: preview,
		})
	}

	return models.CmdStdoutErrorOut{
		OverallNote: err.OverallMessage,
		Errors:      outNotices,
	}
}

func (eb *ErrorBuilder) buildSingleError(refError arch.ReferencedError) models.CmdStdoutErrorOut {
	ref := refError.Reference()

	preview := ""
	if ref.Valid {
		preview, _ = eb.codePrinter.Print(transformRef(ref), codeprinter.CodePrintOpts{
			LineNumbers: true,
			Arrows:      true,
			Highlight:   eb.useColors,
			Mode:        codeprinter.CodePrintModeExtend,
		})
	}

	return models.CmdStdoutErrorOut{
		Errors: []models.StdoutNotice{
			{
				Notice: arch.Notice{
					Message:   refError.Error(),
					Reference: ref,
				},
				Preview: preview,
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