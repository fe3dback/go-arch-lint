package render

import (
	"fmt"
	"path"
	"strings"
)

func (r *Renderer) asciiColorize(color string, value interface{}) (string, error) {
	colorizer := newColorizer(r.colorPrinter)
	out, err := colorizer.colorize(
		color,
		fmt.Sprintf("%s", value),
	)
	if err != nil {
		return "", fmt.Errorf("failed colorize: %s", err)
	}

	return out, nil
}

func (r *Renderer) asciiTrimPrefix(prefix string, value interface{}) string {
	return strings.TrimPrefix(fmt.Sprintf("%s", value), prefix)
}

func (r *Renderer) asciiTrimSuffix(suffix string, value interface{}) string {
	return strings.TrimSuffix(fmt.Sprintf("%s", value), suffix)
}

func (r *Renderer) asciiDefaultValue(def string, value interface{}) string {
	sValue := fmt.Sprintf("%s", value)

	if sValue == "" {
		return def
	}

	return sValue
}

func (r *Renderer) asciiPadLeft(overallLen int, padStr string, value interface{}) string {
	s := fmt.Sprintf("%s", value)

	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func (r *Renderer) asciiPadRight(overallLen int, padStr string, value interface{}) string {
	s := fmt.Sprintf("%s", value)

	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func (r *Renderer) asciiPathDirectory(value interface{}) string {
	return path.Dir(fmt.Sprintf("%s", value))
}
