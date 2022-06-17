package render

import (
	"fmt"
	"path"
	"strconv"
	"strings"
)

func (r *Renderer) asciiColorize(color string, value interface{}) (string, error) {
	colorizer := newColorizer(r.colorPrinter)
	out, err := colorizer.colorize(
		color,
		fmt.Sprintf("%s", value),
	)
	if err != nil {
		return "", fmt.Errorf("failed colorize: %w", err)
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
	s := fmt.Sprintf("%v", value)

	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func (r *Renderer) asciiPadRight(overallLen int, padStr string, value interface{}) string {
	s := fmt.Sprintf("%v", value)

	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func (r *Renderer) asciiLinePrefix(prefix string, value interface{}) string {
	lines := fmt.Sprintf("%s", value)
	result := make([]string, 0)

	for _, line := range strings.Split(lines, "\n") {
		result = append(result, prefix+line)
	}

	return strings.Join(result, "\n")
}

func (r *Renderer) asciiPathDirectory(value interface{}) string {
	return path.Dir(fmt.Sprintf("%v", value))
}

func (r *Renderer) asciiPlus(a, b interface{}) (int, error) {
	iA, err := strconv.Atoi(fmt.Sprintf("%d", a))
	if err != nil {
		return 0, fmt.Errorf("A component of 'plus' is not int: %s", a)
	}

	iB, err := strconv.Atoi(fmt.Sprintf("%d", b))
	if err != nil {
		return 0, fmt.Errorf("B component of 'plus' is not int: %s", b)
	}

	return iA + iB, nil
}

func (r *Renderer) asciiMinus(a, b interface{}) (int, error) {
	iA, err := strconv.Atoi(fmt.Sprintf("%d", a))
	if err != nil {
		return 0, fmt.Errorf("A component of 'minus' is not int: %s", a)
	}

	iB, err := strconv.Atoi(fmt.Sprintf("%d", b))
	if err != nil {
		return 0, fmt.Errorf("B component of 'minus' is not int: %s", b)
	}

	return iA + iB, nil
}

func (r *Renderer) asciiConcat(sources ...interface{}) string {
	result := ""

	for _, source := range sources {
		result += fmt.Sprintf("%v", source)
	}

	return result
}
