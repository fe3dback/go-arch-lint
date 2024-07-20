package ascii

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func (r *ASCII) asciiColorize(color string, value interface{}) string {
	return r.asciiColorizer.Colorize(
		models.ColorName(color),
		fmt.Sprintf("%v", value),
	)
}

func (r *ASCII) asciiTrimPrefix(prefix string, value interface{}) string {
	return strings.TrimPrefix(fmt.Sprintf("%s", value), prefix)
}

func (r *ASCII) asciiTrimSuffix(suffix string, value interface{}) string {
	return strings.TrimSuffix(fmt.Sprintf("%s", value), suffix)
}

func (r *ASCII) asciiDefaultValue(def string, value interface{}) string {
	sValue := fmt.Sprintf("%s", value)

	if sValue == "" {
		return def
	}

	return sValue
}

func (r *ASCII) asciiPadLeft(overallLen int, padStr string, value interface{}) string {
	s := fmt.Sprintf("%v", value)

	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func (r *ASCII) asciiPadRight(overallLen int, padStr string, value interface{}) string {
	s := fmt.Sprintf("%v", value)

	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func (r *ASCII) asciiLinePrefix(prefix string, value interface{}) string {
	lines := fmt.Sprintf("%s", value)
	result := make([]string, 0)

	for _, line := range strings.Split(lines, "\n") {
		result = append(result, prefix+line)
	}

	return strings.Join(result, "\n")
}

func (r *ASCII) asciiPathDirectory(value interface{}) string {
	return path.Dir(fmt.Sprintf("%v", value))
}

func (r *ASCII) asciiPlus(a, b interface{}) (int, error) {
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

func (r *ASCII) asciiMinus(a, b interface{}) (int, error) {
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

func (r *ASCII) asciiConcat(sources ...interface{}) string {
	result := ""

	for _, source := range sources {
		result += fmt.Sprintf("%v", source)
	}

	return result
}
