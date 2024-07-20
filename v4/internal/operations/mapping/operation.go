package mapping

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type Operation struct {
}

func NewOperation() *Operation {
	return &Operation{}
}

func (o *Operation) Mapping(in models.CmdMappingIn) (models.CmdMappingOut, error) {
	// todo:
	return models.CmdMappingOut{
		ProjectDirectory: "todo-ProjectDirectory",
		ModuleName:       "todo-ModuleName",
		MappingGrouped:   nil,
		MappingList: []models.CmdMappingOutList{
			{
				FileName:      "/home/neo/code/fe3dback/go-arch-lint/v4/internal/operations/mapping/operation.go",
				ComponentName: "mapping",
			},
		},
		Scheme: in.Scheme,
	}, nil
}
