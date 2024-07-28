package tests

import (
	"bytes"
	"strings"

	"github.com/kortschak/utter"
)

func Dump(value any) string {
	var golden bytes.Buffer
	utter.Config.ElideType = true
	utter.Config.SortKeys = true
	utter.Fdump(&golden, value)

	result := golden.String()
	result = strings.ReplaceAll(result, "github.com/fe3dback/go-arch-lint/v4/", "")
	result = strings.ReplaceAll(result, "github.com/fe3dback/go-arch-lint/", "")

	return result
}
