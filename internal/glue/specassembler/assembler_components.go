package specassembler

import (
	"fmt"

	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type componentsAssembler struct {
	resolver                *resolver
	allowedImportsAssembler *allowedImportsAssembler
}

func newComponentsAssembler(
	resolver *resolver,
	allowedImportsAssembler *allowedImportsAssembler,
) *componentsAssembler {
	return &componentsAssembler{
		resolver:                resolver,
		allowedImportsAssembler: allowedImportsAssembler,
	}
}

func (m componentsAssembler) assemble(spec *models.ArchSpec, yamlSpec *yaml.YamlSpec) error {
	for yamlName, yamlComponent := range yamlSpec.Components {
		component, err := m.assembleComponent(yamlName, yamlComponent, yamlSpec)
		if err != nil {
			return fmt.Errorf("failed to assemble component: %s", yamlName)
		}

		spec.Components = append(spec.Components, component)
	}

	return nil
}

func (m componentsAssembler) assembleComponent(
	yamlName yaml.YamlComponentName,
	yamlComponent yaml.YamlComponent,
	yamlSpec *yaml.YamlSpec,
) (*models.Component, error) {
	depMeta := yamlSpec.Dependencies[yamlName]

	mayDependOn := make([]models.ComponentName, 0)
	for _, name := range depMeta.MayDependOn {
		mayDependOn = append(mayDependOn, name)
	}

	canUse := make([]models.VendorName, 0)
	for _, name := range depMeta.CanUse {
		canUse = append(canUse, name)
	}

	resolvedPath, err := m.resolver.resolveLocalPath(yamlComponent.LocalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble component path's: %v", err)
	}

	allowedImports, err := m.allowedImportsAssembler.assemble(yamlSpec, mayDependOn, canUse)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble component path's: %v", err)
	}

	return &models.Component{
		Name:           yamlName,
		LocalPathMask:  yamlComponent.LocalPath,
		ResolvedPaths:  resolvedPath,
		MayDependOn:    mayDependOn,
		CanUse:         canUse,
		AllowedImports: allowedImports,
		SpecialFlags: &models.SpecialFlags{
			AllowAllProjectDeps: depMeta.AnyProjectDeps,
			AllowAllVendorDeps:  depMeta.AnyVendorDeps,
		},
	}, nil
}
