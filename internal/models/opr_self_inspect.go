package models

import "github.com/fe3dback/go-arch-lint/internal/models/common"

type (
	CmdSelfInspectIn struct {
		ProjectPath string
		ArchFile    string
	}

	CmdSelfInspectOut struct {
		ModuleName    string                        `json:"ModuleName"`
		RootDirectory string                        `json:"RootDirectory"`
		LinterVersion string                        `json:"LinterVersion"`
		Notices       []CmdSelfInspectOutAnnotation `json:"Notices"`
		Suggestions   []CmdSelfInspectOutAnnotation `json:"Suggestions"`
	}

	CmdSelfInspectOutAnnotation struct {
		Text      string           `json:"Text"`
		Reference common.Reference `json:"Reference"`
	}
)
