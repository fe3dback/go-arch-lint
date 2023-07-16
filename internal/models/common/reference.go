package common

import (
	"fmt"
	"math"
)

type (
	Reference struct {
		Valid    bool
		File     string
		Line     int
		LineFrom int
		LineTo   int
		Column   int
	}
)

func NewReferenceSingleLine(file string, line int, column int) Reference {
	return Reference{Valid: true}.guaranteeValidState(func(r *Reference) {
		r.Valid = true
		r.File = file
		r.Line = line
		r.LineFrom = line
		r.LineTo = line
		r.Column = column
	})
}

func NewReferenceRange(file string, lineFrom, lineMain, lineTo int) Reference {
	return Reference{Valid: true}.guaranteeValidState(func(r *Reference) {
		r.Valid = true
		r.File = file
		r.Line = lineMain
		r.LineFrom = lineFrom
		r.LineTo = lineTo
		r.Column = 0
	})
}

func NewEmptyReference() Reference {
	return Reference{Valid: false}
}

func (r Reference) String() string {
	if !r.Valid {
		return "<unknown_file_ref>"
	}

	return fmt.Sprintf("%s:%d", r.File, r.Line)
}

// ExtendRange will extend from-to range in both ways on growLinesCount lines
// for example initialRange=[2..5], after ExtendRange(1) it will be [1..6]
func (r Reference) ExtendRange(lower, upper int) Reference {
	return r.guaranteeValidState(func(r *Reference) {
		r.LineFrom -= lower
		r.LineTo += upper
	})
}

// ClampWithRealLinesCount allows to clamp lines to real file lines count (upper clamp)
func (r Reference) ClampWithRealLinesCount(linesCount int) Reference {
	return r.guaranteeValidState(func(r *Reference) {
		r.LineFrom = clampInt(r.LineFrom, 1, linesCount)
		r.Line = clampInt(r.Line, 1, linesCount)
		r.LineTo = clampInt(r.LineTo, 1, linesCount)
	})
}

func (r Reference) guaranteeValidState(mutate func(r *Reference)) Reference {
	if !r.Valid {
		return r
	}

	mutate(&r)

	// check lines
	if r.LineFrom > r.LineTo {
		r.LineFrom, r.LineTo = r.LineTo, r.LineFrom
	}

	r.LineFrom = clampInt(r.LineFrom, 1, r.LineTo)
	r.LineTo = clampInt(r.LineTo, r.LineFrom, math.MaxInt32)
	r.Line = clampInt(r.Line, r.LineFrom, r.LineTo)
	r.Column = clampInt(r.Column, 0, math.MaxInt32)

	if r.File == "" {
		r.Valid = false
	}

	return r
}

func clampInt(num, a, b int) int {
	min, max := sortInt(a, b)

	if num < min {
		num = min
	}
	if num > max {
		num = max
	}

	return num
}

func sortInt(a, b int) (min, max int) {
	if a > b {
		return b, a
	}

	return a, b
}
