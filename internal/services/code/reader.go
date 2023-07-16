package code

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

func readFile(fileName string) (content io.Reader, linesCount int) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0
	}

	linesCount, err = lineCounter(file)
	if err != nil {
		return nil, 0
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, 0
	}

	return file, linesCount
}

func readLines(r io.Reader, ref common.Reference) []byte {
	sc := bufio.NewScanner(r)
	currentLine := 0
	var buffer bytes.Buffer

	for sc.Scan() {
		currentLine++

		if currentLine >= ref.LineFrom && currentLine <= ref.LineTo {
			buffer.Write(sc.Bytes())

			if currentLine != ref.LineTo {
				buffer.WriteByte('\n')
			}
		}
	}

	return buffer.Bytes()
}

func highlightContent(filePath string, code []byte) []byte {
	lexer := lexers.Match(filePath)
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
