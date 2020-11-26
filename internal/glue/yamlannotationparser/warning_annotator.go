package yamlannotationparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	sourceMarkerPosition = "^"
)

type AnnotationParser struct {
}

func NewAnnotationParser() *AnnotationParser {
	return &AnnotationParser{}
}

func (a AnnotationParser) Parse(sourceText string) (line, pos int, err error) {
	marker, err := parse(sourceText)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse warning source text: %v", err)
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
//
// Another example without line mark:
// 	   32 |        - 3rd-cobra
// 	   33 |   cmd:    canUse:
// 	                  ^
// 	   36 |       - go-modfile
// 	   37 |
// 	   38 |   a:
func parse(sourceText string) (sourceMarker, error) {
	notValid := sourceMarker{}

	if !strings.Contains(sourceText, sourceMarkerPosition) {
		return notValid, fmt.Errorf("not found position marker")
	}

	marker := sourceMarker{
		sourceLine: 0,
		sourcePos:  0,
	}

	leftOffset := 0
	previousLine := ""

	for _, sourceLine := range strings.Split(sourceText, "\n") {
		if !strings.Contains(sourceLine, sourceMarkerPosition) {
			previousLine = sourceLine
			continue
		}

		// found pos marker line with "^"
		leftOffset = strings.Index(previousLine, `|`) + 1
		marker.sourcePos = strings.Index(sourceLine, sourceMarkerPosition) - leftOffset

		matches := regexp.MustCompile(`^>?\s+(\d+)\s+\|`).FindStringSubmatch(previousLine)
		if len(matches) != 2 {
			return notValid, fmt.Errorf("not found line number in '%s'", previousLine)
		}

		lineNumber, err := strconv.Atoi(matches[1])
		if err != nil {
			return notValid, fmt.Errorf("not found line number in first match")
		}

		marker.sourceLine = lineNumber
	}

	return marker, nil
}
