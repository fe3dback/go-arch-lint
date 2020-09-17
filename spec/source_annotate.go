package spec

import (
	"regexp"
	"strconv"
	"strings"
)

type (
	sourceMarker struct {
		valid bool
		line  int
		pos   int
	}
)

const sourceMarkerLine = ">"
const sourceMarkerPosition = "^"

// Example of sourceText
//	  181 |   game_component:
//	  182 |     mayDependOn:
//	> 183 |       - engine
//	      |         ^
//	  184 |       - engine_entity
//	  185 |       - game_component
//	  186 |       - game_utils
//	  187 |
func parseSourceError(sourceText string) sourceMarker {
	notValid := sourceMarker{valid: false}

	if !strings.Contains(sourceText, sourceMarkerLine) {
		return notValid
	}

	if !strings.Contains(sourceText, sourceMarkerPosition) {
		return notValid
	}

	marker := sourceMarker{
		valid: false,
		line:  0,
		pos:   0,
	}

	lineFound := false
	leftOffset := 0

	for _, sourceLine := range strings.Split(sourceText, "\n") {
		if lineFound {
			// in marker pos line
			// `                ^`

			marker.pos = strings.Index(sourceLine, sourceMarkerPosition) - leftOffset
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

		marker.line = lineNumber
		lineFound = true
		leftOffset = strings.Index(sourceLine, `|`) + 1
	}

	return marker
}
