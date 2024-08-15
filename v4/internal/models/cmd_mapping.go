package models

const (
	CmdMappingSchemeGrouped CmdMappingScheme = "grouped"
	CmdMappingSchemeList    CmdMappingScheme = "list"
)

var CmdMappingSchemesValues = []string{
	CmdMappingSchemeList,
	CmdMappingSchemeGrouped,
}

type (
	CmdMappingScheme = string

	CmdMappingIn struct {
		Scheme CmdMappingScheme
	}

	CmdMappingOut struct {
		ProjectDirectory string                 `json:"ProjectDirectory"`
		ModuleName       string                 `json:"ModuleName"`
		MappingGrouped   []CmdMappingOutGrouped `json:"MappingGrouped"`
		MappingList      []CmdMappingOutList    `json:"MappingList"`
		Scheme           CmdMappingScheme       `json:"-"`
	}

	CmdMappingOutGrouped struct {
		ComponentName string   `json:"ComponentName"`
		Packages      []string `json:"Packages"`
	}

	CmdMappingOutList struct {
		Package       string `json:"Package"`
		ComponentName string `json:"ComponentName"`
	}
)
