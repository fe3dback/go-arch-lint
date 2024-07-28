package mapping

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

// todo: remove
type ConfRndTmp interface {
	FetchSpec(path models.PathAbsolute) (models.Config, error)
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
	conf, err := o.ConfRndTmp.FetchSpec("/home/neo/code/fe3dback/go-arch-lint/v4/.go-arch-lint.yml")
	if err != nil {
		return models.CmdMappingOut{}, fmt.Errorf("failed read config: %w", err)
	}

	fmt.Printf("DS global: value=%v at ref=[valid=%v]\n", conf.Settings.DeepScan.Value, conf.Settings.DeepScan.Ref.Valid)

	conf.Dependencies.Map.Each(func(name models.ComponentName, rules models.ConfigComponentDependencies, _ models.Reference) {
		fmt.Printf("- DS '%s': defined=%v, value=%v at ref=[valid=%v]\n",
			name,
			rules.DeepScan.Defined,
			rules.DeepScan.Value.Value,
			rules.DeepScan.Value.Ref.Valid,
		)
	})

	_ = conf

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
