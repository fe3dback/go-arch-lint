package models

type (
	ProjectInfo struct {
		Directory  PathAbsolute
		ConfigPath PathAbsolute
		Module     GoModule
	}
)
