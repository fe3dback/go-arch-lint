package warnparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const sourceMarkerLine = ">"
const sourceMarkerPosition = "^"

type WarningSourceParser struct {
}

func NewWarningSourceParser() *WarningSourceParser {
	return &WarningSourceParser{}
}

func (a WarningSourceParser) Parse(sourceText string) (line, pos int, err error) {
	marker := parseSourceWarning(sourceText)
	if !marker.valid {
		return 0, 0, fmt.Errorf("failed to parse warning source text")
	}

	return marker.sourceLine, marker.sourcePos, nil
}

// Example of sourceText
//	  181 |   game_component:
//	  182 |     mayDependOn:
//	> 183 |       - engine
//	      |         ^
//	  184 |       - engine_entity
//	  185 |       - game_component
//	  186 |       - game_utils
//	  187 |
func parseSourceWarning(sourceText string) sourceMarker {
	notValid := sourceMarker{valid: false}

	if !strings.Contains(sourceText, sourceMarkerLine) {
		return notValid
	}

	if !strings.Contains(sourceText, sourceMarkerPosition) {
		return notValid
	}

	marker := sourceMarker{
		valid:      false,
		sourceLine: 0,
		sourcePos:  0,
	}

	lineFound := false
	leftOffset := 0

	for _, sourceLine := range strings.Split(sourceText, "\n") {
		if lineFound {
			// in marker pos line
			// `                ^`

			marker.sourcePos = strings.Index(sourceLine, sourceMarkerPosition) - leftOffset
			marker.valid = true
			break
		}

		if !strings.Contains(sourceLine, sourceMarkerLine) {
			continue
		}

		// in marker line
		// `> 183 |       - engine`

		matches := regexp.MustCompile(`^>\s+(\d+)\s+\|`).FindStringSubmatch(sourceLine)
		if len(matches) != 2 {
			return notValid
		}

		lineNumber, err := strconv.Atoi(matches[1])
		if err != nil {
			return notValid
		}

		marker.sourceLine = lineNumber
		lineFound = true
		leftOffset = strings.Index(sourceLine, `|`) + 1
	}

	return marker
}
