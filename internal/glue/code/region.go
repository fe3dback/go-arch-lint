package code

type (
	codeRegion struct {
		lineFirst int
		lineMain  int
		lineLast  int
	}
)

func calculateCodeRegion(line int, regionHeight int, maxLines int) codeRegion {
	lineFirst := line
	lineMain := line
	lineLast := line + (regionHeight - 1)

	if lineFirst < 1 {
		lineFirst = 1
	}

	if lineLast > maxLines {
		lineLast = maxLines
	}

	return codeRegion{
		lineFirst: lineFirst,
		lineMain:  lineMain,
		lineLast:  lineLast,
	}
}
