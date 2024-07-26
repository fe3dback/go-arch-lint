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
		ComponentName string
		FileNames     []string
	}

	CmdMappingOutList struct {
		FileName      string
		ComponentName string
	}
)
