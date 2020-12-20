package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	componentsAssembler struct {
		resolver                *resolver
		allowedImportsAssembler *allowedImportsAssembler
	}
)

func newComponentsAssembler(
	resolver *resolver,
	allowedImportsAssembler *allowedImportsAssembler,
) *componentsAssembler {
	return &componentsAssembler{
		resolver:                resolver,
		allowedImportsAssembler: allowedImportsAssembler,
	}
}

func (m componentsAssembler) assemble(spec *speca.Spec, document arch.Document) error {
	for yamlName, yamlComponent := range document.Components().Map() {
		component, err := m.assembleComponent(yamlName, yamlComponent, document)
		if err != nil {
			return fmt.Errorf("failed to assemble component '%s': %w", yamlName, err)
		}

		spec.Components = append(spec.Components, component)
	}

	return nil
}

func (m componentsAssembler) assembleComponent(
	yamlName string,
	yamlComponent arch.Component,
	yamlDocument arch.Document,
) (speca.Component, error) {
	depMeta, hasDeps := yamlDocument.Dependencies().Map()[yamlName]

	mayDependOn := make([]speca.ReferableString, 0)
	canUse := make([]speca.ReferableString, 0)

	if hasDeps {
		for _, componentName := range depMeta.MayDependOn() {
			mayDependOn = append(mayDependOn, componentName)
		}

		for _, vendorName := range depMeta.CanUse() {
			canUse = append(canUse, vendorName)
		}
	}

	// component path in
	tmpResolvedPath, err := m.resolver.resolveLocalPath(yamlComponent.LocalPath().Value())
	if err != nil {
		return speca.Component{}, fmt.Errorf("failed to assemble component path's: %v", err)
	}
	resolvedPaths := wrapPaths(
		yamlComponent.LocalPath().Reference(),
		tmpResolvedPath,
	)

	// deps import
	tmpAllowedImports, err := m.allowedImportsAssembler.assemble(
		yamlDocument,
		unwrapStrings(mayDependOn),
		unwrapStrings(canUse),
	)
	if err != nil {
		return speca.Component{}, fmt.Errorf("failed to assemble component path's: %v", err)
	}
	allowedImports := wrapPaths(
		yamlComponent.Reference(),
		tmpAllowedImports,
	)

	var specialFlags speca.SpecialFlags
	if !hasDeps {
		specialFlags = speca.SpecialFlags{
			AllowAllProjectDeps: speca.NewReferableBool(false, yamlComponent.Reference()),
			AllowAllVendorDeps:  speca.NewReferableBool(false, yamlComponent.Reference()),
		}
	} else {
		specialFlags = speca.SpecialFlags{
			AllowAllProjectDeps: depMeta.AnyProjectDeps(),
			AllowAllVendorDeps:  depMeta.AnyVendorDeps(),
		}
	}

	return speca.Component{
		Name: speca.NewReferableString(
			yamlName,
			yamlComponent.Reference(),
		),
		LocalPathMask:  yamlComponent.LocalPath(),
		ResolvedPaths:  resolvedPaths,
		MayDependOn:    mayDependOn,
		CanUse:         canUse,
		AllowedImports: allowedImports,
		SpecialFlags:   specialFlags,
	}, nil
}
