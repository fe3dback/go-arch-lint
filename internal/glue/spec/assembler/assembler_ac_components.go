package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	componentsAssembler struct {
		resolver                *resolver
		allowedImportsAssembler *allowedImportsAssembler
		provideYamlRef          provideYamlRef
	}
)

func newComponentsAssembler(
	resolver *resolver,
	allowedImportsAssembler *allowedImportsAssembler,
	provideYamlRef provideYamlRef,
) *componentsAssembler {
	return &componentsAssembler{
		resolver:                resolver,
		allowedImportsAssembler: allowedImportsAssembler,
		provideYamlRef:          provideYamlRef,
	}
}

func (m componentsAssembler) assemble(spec *speca.Spec, yamlSpec *spec.ArchV1Document) error {
	for yamlName, yamlComponent := range yamlSpec.V1Components {
		component, err := m.assembleComponent(yamlName, yamlComponent, yamlSpec)
		if err != nil {
			return fmt.Errorf("failed to assemble component: %s", yamlName)
		}

		spec.Components = append(spec.Components, component)
	}

	return nil
}

func (m componentsAssembler) assembleComponent(
	yamlName string,
	yamlComponent spec.ArchV1Component,
	yamlSpec *spec.ArchV1Document,
) (speca.Component, error) {
	depMeta := yamlSpec.V1Dependencies[yamlName]

	// components
	mayDependOn := make([]speca.ReferableString, 0)
	for index, name := range depMeta.V1MayDependOn {
		mayDependOn = append(mayDependOn, speca.NewReferableString(
			name,
			m.provideYamlRef(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", yamlName, index)),
		))
	}

	// vendors
	canUse := make([]speca.ReferableString, 0)
	for index, name := range depMeta.V1CanUse {
		canUse = append(canUse, speca.NewReferableString(
			name,
			m.provideYamlRef(fmt.Sprintf("$.deps.%s.canUse[%d]", yamlName, index)),
		))
	}

	// component path in
	tmpResolvedPath, err := m.resolver.resolveLocalPath(yamlComponent.V1LocalPath)
	if err != nil {
		return speca.Component{}, fmt.Errorf("failed to assemble component path's: %v", err)
	}
	resolvedPaths := wrapPaths(
		m.provideYamlRef(fmt.Sprintf("$.components.%s.in", yamlName)),
		tmpResolvedPath,
	)

	// deps import
	tmpAllowedImports, err := m.allowedImportsAssembler.assemble(
		yamlSpec,
		unwrapStrings(mayDependOn),
		unwrapStrings(canUse),
	)
	if err != nil {
		return speca.Component{}, fmt.Errorf("failed to assemble component path's: %v", err)
	}
	allowedImports := wrapPaths(
		m.provideYamlRef(fmt.Sprintf("$.components.%s", yamlName)),
		tmpAllowedImports,
	)

	return speca.Component{
		Name: speca.NewReferableString(
			yamlName,
			m.provideYamlRef(fmt.Sprintf("$.components.%s", yamlName)),
		),
		LocalPathMask: speca.NewReferableString(
			yamlComponent.V1LocalPath,
			m.provideYamlRef(fmt.Sprintf("$.components.%s.in", yamlName)),
		),
		ResolvedPaths:  resolvedPaths,
		MayDependOn:    mayDependOn,
		CanUse:         canUse,
		AllowedImports: allowedImports,
		SpecialFlags: speca.SpecialFlags{
			AllowAllProjectDeps: speca.NewReferableBool(
				depMeta.V1AnyProjectDeps,
				m.provideYamlRef(fmt.Sprintf("$.deps.%s.anyProjectDeps", yamlName)),
			),
			AllowAllVendorDeps: speca.NewReferableBool(
				depMeta.V1AnyVendorDeps,
				m.provideYamlRef(fmt.Sprintf("$.deps.%s.anyVendorDeps", yamlName)),
			),
		},
	}, nil
}
