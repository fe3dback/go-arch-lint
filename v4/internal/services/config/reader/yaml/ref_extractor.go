package yaml

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

var lineRegexp = regexp.MustCompile(`>\s+(\d+)\s\|`)
var lineClearRegexp = regexp.MustCompile(`[\s\n\r\t]`)

// this will parse custom yaml error with source annotation
// and extract File/Line/Pos info to common Reference struct
//
// example:
//
//	]   6 |   imports:
//	]   7 |     strictMode: false
//	]   8 |     allowAnyVendorImports: true
//	]>  9 |   structTags:
//	]         ^
//	]  10 |     allowed: true
//	]  11 |     # allowed: true                # all tags (default)
//	]  12 |     # allowed: false               # no tags
//	]  13 |
func extractReferenceFromError(tCtx TransformContext, errLines []string) models.Reference {
	return models.NewReference(
		tCtx.file,
		findLineNumber(errLines),
		findOffset(errLines),
		"",
	)
}
func findLineNumber(errLines []string) int {
	isSimpleCase := false
	simpleLine := ""

	for _, line := range errLines {
		matched := lineRegexp.FindStringSubmatch(line)
		if len(matched) > 1 {
			isSimpleCase = true
			simpleLine = line
			break
		}
	}

	if isSimpleCase {
		return findLineNumberSimple(simpleLine)
	}

	return findLineNumberComplex(errLines)
}

// ]>  9 |   structTags:
// ]> 24 |   go-common:           { in2: golang.org/x/sync/errgroup }
func findLineNumberSimple(line string) int {
	number := lineRegexp.FindStringSubmatch(line)[1]
	lineNumber, err := strconv.ParseInt(number, 10, 32)
	if err != nil {
		return 0
	}

	return int(lineNumber)
}

// ]  21 |     - "^.*\/test\/.*$":
// ]       ^
// ]  23 | vendors:
func findLineNumberComplex(errLines []string) int {
	prevLine := ""

	for _, line := range errLines {
		cleared := lineClearRegexp.ReplaceAllString(line, "")
		if cleared == "^" {
			return findLineNumberSimple("> "+prevLine) + 1
		}

		prevLine = line
	}

	return 0
}

// ]>  9 |   structTags:
// ]         ^
func findOffset(errLines []string) int {
	for _, line := range errLines {
		cleared := lineClearRegexp.ReplaceAllString(line, "")
		if cleared != "^" {
			continue
		}

		return strings.IndexRune(line, '^') - findCodeOffset(errLines)
	}

	return 0
}

// ]  13 |
// ]>  9 |   structTags:
// ]         ^
// ]  10 |     allowed: true
// ]  11 |     # allowed: true                # all tags (default)
// ]  12 |     # allowed: false               # no tags
func findCodeOffset(errLines []string) int {
	if len(errLines) <= 0 {
		return 0
	}

	firstLine := errLines[0]
	return strings.IndexRune(firstLine, '|') + 1
}
