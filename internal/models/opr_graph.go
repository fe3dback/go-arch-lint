package models

const (
	GraphTypeFlow GraphType = "flow"
	GraphTypeDI   GraphType = "di"
)

var GraphTypesValues = []string{
	GraphTypeFlow,
	GraphTypeDI,
}

type (
	GraphType = string

	CmdGraphIn struct {
		ProjectPath    string
		ArchFile       string
		Type           GraphType
		OutFile        string
		Focus          string
		IncludeVendors bool
	}

	CmdGraphOut struct {
		ProjectDirectory string `json:"ProjectDirectory"`
		ModuleName       string `json:"ModuleName"`
		OutFile          string `json:"OutFile"`
	}
)
