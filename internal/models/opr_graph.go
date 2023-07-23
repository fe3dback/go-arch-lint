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
		ExportD2       bool
		OutputType     OutputType
	}

	CmdGraphOut struct {
		ProjectDirectory string `json:"ProjectDirectory"`
		ModuleName       string `json:"ModuleName"`
		OutFile          string `json:"OutFile"`
		D2Definitions    string `json:"D2Definitions"`
		ExportD2         bool   `json:"-"`
	}
)
