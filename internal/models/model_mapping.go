package models

const (
	MappingSchemeGrouped MappingScheme = "grouped"
	MappingSchemeList    MappingScheme = "list"
)

type (
	MappingScheme = string

	Mapping struct {
		ProjectDirectory string           `json:"ProjectDirectory"`
		ModuleName       string           `json:"ModuleName"`
		MappingGrouped   []MappingGrouped `json:"MappingGrouped"`
		MappingList      []MappingList    `json:"MappingList"`
		Scheme           MappingScheme    `json:"-"`
	}

	MappingGrouped struct {
		ComponentName string
		FileNames     []string
	}

	MappingList struct {
		FileName      string
		ComponentName string
	}
)
