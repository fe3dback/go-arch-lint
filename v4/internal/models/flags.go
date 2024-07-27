package models

const (
	FlagOutputType                  = "output-type"
	FlagOutputUseAsciiColors        = "output-color"
	FlagOutputTypeJSON              = "json"
	FlagOutputJSONWithoutFormatting = "output-json-one-line"
	FlagSkipMissUsages              = "skip-missuse"
	FlagProjectPath                 = "project-path"
	FlagArchConfigRelativePath      = "arch-file"
)

const (
	FlagCategoryGlobal  = "global:"
	FlagCategoryCommand = "this command:"
)

const (
	DefaultProjectPath   = "./"
	DefaultArchFileName  = ".go-arch-lint.yml"
	DefaultGoModFileName = "go.mod"
)
