package xstdout

import (
	"errors"

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
	noticeErr := &models.ErrorWithNotices{}
	if errors.As(err, &noticeErr) {
		return eb.buildNoticesError(noticeErr)
	}

	refError := models.ReferencedError{}
	if errors.As(err, &refError) {
		return eb.buildSingleError(refError)
	}

	return models.CmdStdoutErrorOut{
		Errors: []models.StdoutNotice{
			{
				Notice: models.Notice{
					Message:   err.Error(),
					Reference: models.NewInvalidReference(),
				},
			},
		},
	}
}

func (eb *ErrorBuilder) buildNoticesError(err *models.ErrorWithNotices) models.CmdStdoutErrorOut {
	outNotices := make([]models.StdoutNotice, 0, len(err.Notices))

	mode := models.CodePrintModeExtend
	if len(err.Notices) > 10 {
		mode = models.CodePrintModeOneLine
	}

	for _, notice := range err.Notices {
		preview := ""
		if notice.Reference.Valid {
			preview, _ = eb.codePrinter.Print(notice.Reference, models.CodePrintOpts{
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

func (eb *ErrorBuilder) buildSingleError(refError models.ReferencedError) models.CmdStdoutErrorOut {
	ref := refError.Reference()

	preview := ""
	if ref.Valid {
		preview, _ = eb.codePrinter.Print(ref, models.CodePrintOpts{
			LineNumbers: true,
			Arrows:      true,
			Highlight:   eb.useColors,
			Mode:        models.CodePrintModeExtend,
		})
	}

	return models.CmdStdoutErrorOut{
		Errors: []models.StdoutNotice{
			{
				Notice: models.Notice{
					Message:   refError.Error(),
					Reference: ref,
				},
				Preview: preview,
			},
		},
	}
}
