package code

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Render struct {
	printer ColorPrinter
}

func NewRender(printer ColorPrinter) *Render {
	return &Render{
		printer: printer,
	}
}

func (r *Render) SourceCode(ref models.Reference, height int, highlight bool) []byte {
	code, region := r.readCode(ref, height, highlight)
	return r.annotate(code, ref, region)
}

func (r *Render) readCode(ref models.Reference, height int, highlight bool) ([]byte, codeRegion) {
	if !ref.Valid {
		return []byte{}, codeRegion{}
	}

	rawCode, region := r.readRaw(ref, height)
	if !highlight {
		return rawCode, region
	}

	return highlightRawCode(ref, rawCode), region
}

func (r *Render) readRaw(ref models.Reference, height int) ([]byte, codeRegion) {
	if !ref.Valid {
		return []byte{}, codeRegion{}
	}

	file, err := os.Open(ref.File)
	if err != nil {
		return []byte{}, codeRegion{}
	}

	linesCount, err := lineCounter(file)
	if err != nil {
		return []byte{}, codeRegion{}
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return []byte{}, codeRegion{}
	}

	region := calculateCodeRegion(ref.Line, height, linesCount)
	return readLines(file, region), region
}

func highlightRawCode(ref models.Reference, code []byte) []byte {
	lexer := lexers.Match(ref.File)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Trac
	formatter := formatters.TTY8

	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		return []byte{}
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return []byte{}
	}

	return buf.Bytes()
}

func readLines(r io.Reader, region codeRegion) []byte {
	sc := bufio.NewScanner(r)
	currentLine := 0
	var buffer bytes.Buffer

	for sc.Scan() {
		currentLine++

		if currentLine >= region.lineFirst && currentLine <= region.lineLast {
			buffer.Write(sc.Bytes())

			if currentLine != region.lineLast {
				buffer.WriteByte('\n')
			}
		}
	}

	return buffer.Bytes()
}

func (r *Render) annotate(
	code []byte,
	ref models.Reference,
	region codeRegion,
) []byte {
	buf := bytes.NewBuffer(code)
	sc := bufio.NewScanner(buf)
	currentLine := region.lineFirst

	var resultBuffer bytes.Buffer
	for sc.Scan() {
		prefixLine := r.printer.Gray(fmt.Sprintf("%4d |", currentLine))
		prefixEmpty := r.printer.Gray("     |")
		resultBuffer.WriteString(fmt.Sprintf("%s %s\n", prefixLine, sc.Bytes()))

		if currentLine == region.lineMain && ref.Valid {
			spaces := strings.Repeat(" ", int(math.Max(0, float64(ref.Offset-1))))
			resultBuffer.WriteString(fmt.Sprintf("%s %s^\n", prefixEmpty, spaces))
		}

		currentLine++
	}

	return resultBuffer.Bytes()
}
