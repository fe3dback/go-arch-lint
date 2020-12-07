package models

type (
	FlagsCheck struct {
		ProjectDirectory string
		GoArchFilePath   string
		GoModFilePath    string
		ModuleName       string
		MaxWarnings      int
	}
)
