package ascii

import (
	"strings"
)

func preprocessRawASCIITemplate(tpl string) string {
	lines := strings.Split(tpl, "\n")
	buffer := make([]string, 0, len(lines))

	for _, line := range lines {
		processedLine := preprocessTemplateLine(line)
		if strings.TrimSpace(processedLine) == "" {
			continue
		}

		buffer = append(buffer, processedLine)
	}

	return strings.Join(buffer, "\n")
}

func preprocessTemplateLine(row string) string {
	return strings.TrimLeft(row, "\t")
}
