package mapping

import (
	"fmt"
	"sort"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Operation struct {
	specFetcher specFetcher
}

func NewOperation(specFetcher specFetcher) *Operation {
	return &Operation{
		specFetcher: specFetcher,
	}
}

func (o *Operation) Mapping(in models.CmdMappingIn) (models.CmdMappingOut, error) {
	spec, err := o.specFetcher.FetchSpec()
	if err != nil {
		return models.CmdMappingOut{}, fmt.Errorf("failed read config: %w", err)
	}

	return models.CmdMappingOut{
		ProjectDirectory: string(spec.Project.Directory),
		ModuleName:       string(spec.Project.Module),
		MappingGrouped:   o.buildGrouped(spec),
		MappingList:      o.buildList(spec),
		Scheme:           in.Scheme,
	}, nil
}

func (o *Operation) buildGrouped(spec models.Spec) []models.CmdMappingOutGrouped {
	list := make([]models.CmdMappingOutGrouped, 0, len(spec.Components))

	for _, component := range spec.Components {
		group := models.CmdMappingOutGrouped{
			ComponentName: string(component.Name.Value),
			Packages:      make([]string, 0, len(component.OwnedPackages)),
		}

		for _, pkg := range component.OwnedPackages {
			group.Packages = append(group.Packages, string(pkg.PathAbs))
		}

		list = append(list, group)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ComponentName <= list[j].ComponentName
	})

	return list
}

func (o *Operation) buildList(spec models.Spec) []models.CmdMappingOutList {
	list := make([]models.CmdMappingOutList, 0, 128)

	for _, component := range spec.Components {
		for _, ownedPackage := range component.OwnedPackages {
			list = append(list, models.CmdMappingOutList{
				Package:       string(ownedPackage.PathAbs),
				ComponentName: string(component.Name.Value),
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ComponentName <= list[j].ComponentName
	})

	return list
}
