package code

import (
	"math"
)

type (
	codeRegion struct {
		lineFirst int
		lineMain  int
		lineLast  int
	}
)

func calculateCodeRegion(line int, regionHeight int, maxLines int) codeRegion {
	if line < 1 {
		line = 1
	}
	if line > maxLines {
		line = maxLines
	}

	isEven := regionHeight%2 == 0
	var topHeight int

	if isEven {
		topHeight = regionHeight / 2 // 10 -> 5
	} else {
		topHeight = (regionHeight + 1) / 2 // 9 -> 5
	}

	firstLine := line - topHeight
	if firstLine < 1 {
		firstLine = 1
	}

	lastLine := firstLine + regionHeight
	if lastLine > maxLines {
		lastLine = maxLines
	}

	if (lastLine - firstLine) < regionHeight {
		firstLine = int(math.Max(1, float64(lastLine-regionHeight)))
	}

	return codeRegion{
		lineFirst: firstLine,
		lineMain:  line,
		lineLast:  lastLine,
	}
}
