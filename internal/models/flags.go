package models

const (
	DefaultProjectPath   = "./"
	DefaultArchFileName  = ".go-arch-lint.yml"
	DefaultGoModFileName = "go.mod"
)

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
