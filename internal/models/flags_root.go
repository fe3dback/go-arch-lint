package models

type (
	OutputType = string

	FlagsRoot struct {
		UseColors         bool
		OutputType        OutputType
		OutputJsonOneLine bool
	}
)

var OutputTypeVariantsConst = []string{
	OutputTypeASCII,
	OutputTypeJSON,
}

const (
	OutputTypeDefault OutputType = "default"
	OutputTypeASCII   OutputType = "ascii"
	OutputTypeJSON    OutputType = "json"
)
