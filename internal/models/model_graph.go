package models

const (
	GraphTypeFlow GraphType = "flow"
	GraphTypeDI   GraphType = "di"
)

type (
	GraphType = string

	Graph struct {
		ProjectDirectory string `json:"ProjectDirectory"`
		ModuleName       string `json:"ModuleName"`
		OutFile          string `json:"OutFile"`
	}
)
