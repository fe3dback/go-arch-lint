package codeprinter

import (
	"fmt"
	"os"
	"strings"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type ExtractorRaw struct{}

func NewExtractorRaw() *ExtractorRaw {
	return &ExtractorRaw{}
}

func (e *ExtractorRaw) ExtractLines(file arch.PathAbsolute, from int, to int) ([]string, error) {
	data, err := os.ReadFile(string(file))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	return safeTakeLines(lines, from, to), nil
}
