package codeprinter

func safeTakeLines(in []string, from, to int) []string {
	lines := len(in)
	indStart := clamp(from, 0, lines)
	indEnd := clamp(to, 0, lines)

	return in[indStart-1 : indEnd]
}

func clamp(v, lb, ub int) int {
	if v < lb {
		v = lb
	}

	if v > ub {
		v = ub
	}

	return v
}
