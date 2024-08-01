package mapping

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

// todo: remove
type ConfRndTmp interface {
	FetchSpec(path models.PathAbsolute) (models.Spec, error)
}

type Operation struct {
	ConfRndTmp ConfRndTmp
}

func NewOperation(ConfRndTmp ConfRndTmp) *Operation {
	return &Operation{
		ConfRndTmp: ConfRndTmp,
	}
}

func (o *Operation) Mapping(in models.CmdMappingIn) (models.CmdMappingOut, error) {
	// todo: module info fetcher
	//  - this will found dir/conf file/module name and other base info

	// todo: remove
	spec, err := o.ConfRndTmp.FetchSpec("/home/neo/code/fe3dback/go-arch-lint/v4/.go-arch-lint.yml")
	if err != nil {
		return models.CmdMappingOut{}, fmt.Errorf("failed read config: %w", err)
	}

	_ = spec

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
