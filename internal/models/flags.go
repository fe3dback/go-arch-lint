package models

const (
	OutputTypeDefault OutputType = "default"
	OutputTypeASCII   OutputType = "ascii"
	OutputTypeJSON    OutputType = "json"
)

var OutputTypeValues = []string{
	OutputTypeASCII,
	OutputTypeJSON,
}

type (
	OutputType = string

	FlagsRoot struct {
		UseColors         bool
		OutputType        OutputType
		OutputJsonOneLine bool
	}
)
