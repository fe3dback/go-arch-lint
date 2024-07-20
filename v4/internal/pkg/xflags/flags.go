package xflags

import (
	"fmt"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
)

func CreateEnumFlag(name string, aliases []string, usage string, values []string, defaultValue string, category string) cli.Flag {
	if !slices.Contains(values, defaultValue) {
		panic(fmt.Sprintf("unknown flag %s default value '%v' (possible is '%#v')", name, defaultValue, values))
	}

	commaList := "[" + strings.Join(values, ", ") + "]"

	return &cli.StringFlag{
		Name:     name,
		Aliases:  aliases,
		Category: category,
		Usage:    fmt.Sprintf("%s. Enum:%s", usage, commaList),
		Value:    defaultValue,
		Action: func(context *cli.Context, value string) error {
			if !slices.Contains(values, value) {
				return fmt.Errorf("command flag '%s' has unknown value '%s'. Expected one of %s", name, value, commaList)
			}

			return nil
		},
	}
}
