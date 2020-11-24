package render

import (
	"strings"
)

func preprocessRawAsciiTemplate(tpl string) string {
	lines := strings.Split(tpl, "\n")
	buffer := make([]string, 0, len(lines))

	for _, line := range lines {
		processedLine := preprocessTemplateLine(line)
		if processedLine == "" {
			continue
		}

		buffer = append(buffer, processedLine)
	}

	return strings.Join(buffer, "\n")
}

func preprocessTemplateLine(row string) string {
	return strings.ReplaceAll(
		strings.TrimPrefix(row, "\t"),
		"\t",
		"  ",
	)
}
