package models

const (
	ProjectInfoDefaultArchFileName  = ".go-arch-lint.yml"
	ProjectInfoDefaultGoModFileName = "go.mod"
)

type (
	ProjectInfo struct {
		Directory      string
		GoArchFilePath string
		GoModFilePath  string
		ModuleName     string
	}
)
