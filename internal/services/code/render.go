package code

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type Render struct {
	printer colorPrinter
}

type annotateOpts struct {
	code              []byte
	ref               common.Reference
	showColumnPointer bool
}

func NewRender(printer colorPrinter) *Render {
	return &Render{
		printer: printer,
	}
}

func (r *Render) SourceCode(ref common.Reference, highlight bool, showPointer bool) []byte {
	opts := r.fetch(ref, highlight)
	opts.showColumnPointer = showPointer

	return r.annotate(opts)
}

func (r *Render) fetch(ref common.Reference, highlight bool) annotateOpts {
	if !ref.Valid {
		return annotateOpts{}
	}

	file, linesCount := readFile(ref.File)
	if file == nil {
		return annotateOpts{}
	}

	ref = ref.ClampWithRealLinesCount(linesCount)
	content := readLines(file, ref)

	if highlight {
		content = highlightContent(ref.File, content)
	}

	return annotateOpts{code: content, ref: ref}
}

func (r *Render) annotate(opt annotateOpts) []byte {
	buf := bytes.NewBuffer(opt.code)
	sc := bufio.NewScanner(buf)
	currentLine := opt.ref.LineFrom

	var resultBuffer bytes.Buffer
	for sc.Scan() {
		prefixLine := r.printer.Gray(fmt.Sprintf("%4d |", currentLine))
		prefixEmpty := r.printer.Gray("        ")

		// add line pointer
		if currentLine == opt.ref.Line {
			prefixLine = fmt.Sprintf("> %s", prefixLine)
		} else {
			prefixLine = fmt.Sprintf("  %s", prefixLine)
		}

		// draw line
		resultBuffer.WriteString(fmt.Sprintf("%s %s\n",
			prefixLine,
			r.replaceTabsToSpaces(sc.Bytes()),
		))

		// add offset pointer
		if opt.showColumnPointer {
			if currentLine == opt.ref.Line {
				spaces := strings.Repeat(" ", int(math.Max(0, float64(opt.ref.Column-1))))
				resultBuffer.WriteString(fmt.Sprintf("%s %s^\n", prefixEmpty, spaces))
			}
		}

		currentLine++
	}

	return resultBuffer.Bytes()
}

func (r *Render) replaceTabsToSpaces(src []byte) []byte {
	return []byte(strings.ReplaceAll(string(src), "\t", "  "))
}
