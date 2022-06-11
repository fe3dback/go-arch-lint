package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	componentsAssembler struct {
		resolver                       *resolver
		allowedProjectImportsAssembler *allowedProjectImportsAssembler
		allowedVendorImportsAssembler  *allowedVendorImportsAssembler
	}
)

func newComponentsAssembler(
	resolver *resolver,
	allowedProjectImportsAssembler *allowedProjectImportsAssembler,
	allowedVendorImportsAssembler *allowedVendorImportsAssembler,
) *componentsAssembler {
	return &componentsAssembler{
		resolver:                       resolver,
		allowedProjectImportsAssembler: allowedProjectImportsAssembler,
		allowedVendorImportsAssembler:  allowedVendorImportsAssembler,
	}
}

func (m *componentsAssembler) assemble(spec *speca.Spec, document arch.Document) error {
	for yamlName, yamlComponent := range document.Components().Map() {
		component, err := m.assembleComponent(yamlName, yamlComponent, document)
		if err != nil {
			return fmt.Errorf("failed to assemble component '%s': %w", yamlName, err)
		}

		spec.Components = append(spec.Components, component)
	}

	return nil
}

func (m *componentsAssembler) assembleComponent(
	yamlName string,
	yamlComponent arch.Component,
	yamlDocument arch.Document,
) (speca.Component, error) {
	depMeta, hasDeps := yamlDocument.Dependencies().Map()[yamlName]

	mayDependOn := make([]speca.Referable[string], 0)
	canUse := make([]speca.Referable[string], 0)

	if hasDeps {
		mayDependOn = append(mayDependOn, depMeta.MayDependOn()...)
		canUse = append(canUse, depMeta.CanUse()...)
	}

	cmp := speca.Component{
		Name:        speca.NewReferable(yamlName, yamlComponent.Reference()),
		MayDependOn: mayDependOn,
		CanUse:      canUse,
	}

	type enricher func() error
	enrichers := []enricher{
		func() error { return m.enrichWithFlags(&cmp, yamlComponent, hasDeps, depMeta) },
		func() error { return m.enrichWithResolvedPaths(&cmp, yamlName, yamlComponent) },
		func() error { return m.enrichWithProjectImports(&cmp, yamlComponent, yamlDocument, mayDependOn) },
		func() error { return m.enrichWithVendorGlobs(&cmp, yamlDocument, canUse) },
	}

	for _, enrich := range enrichers {
		err := enrich()
		if err != nil {
			return speca.Component{}, fmt.Errorf("failed assemble component '%s', enrich '%T' err: %w",
				yamlName,
				enrich,
				err,
			)
		}
	}

	return cmp, nil
}

func (m *componentsAssembler) enrichWithFlags(
	cmp *speca.Component,
	yamlComponent arch.Component,
	hasDeps bool,
	depMeta arch.DependencyRule,
) error {
	if hasDeps {
		cmp.SpecialFlags = speca.SpecialFlags{
			AllowAllProjectDeps: depMeta.AnyProjectDeps(),
			AllowAllVendorDeps:  depMeta.AnyVendorDeps(),
		}
		return nil
	}

	cmp.SpecialFlags = speca.SpecialFlags{
		AllowAllProjectDeps: speca.NewReferable(false, yamlComponent.Reference()),
		AllowAllVendorDeps:  speca.NewReferable(false, yamlComponent.Reference()),
	}

	return nil
}

func (m *componentsAssembler) enrichWithResolvedPaths(
	cmp *speca.Component,
	yamlName string,
	yamlComponent arch.Component,
) error {
	resolvedPaths := make([]speca.Referable[models.ResolvedPath], 0)

	for _, componentIn := range yamlComponent.RelativePaths() {
		tmpResolvedPath, err := m.resolver.resolveLocalGlobPath(string(componentIn.Value()))
		if err != nil {
			return fmt.Errorf("failed to assemble component '%s' path '%s': %w",
				yamlName,
				componentIn.Value(),
				err,
			)
		}

		wrappedPaths := wrap(componentIn.Reference(), tmpResolvedPath)
		resolvedPaths = append(resolvedPaths, wrappedPaths...)
	}

	cmp.ResolvedPaths = resolvedPaths
	return nil
}

func (m *componentsAssembler) enrichWithProjectImports(
	cmp *speca.Component,
	yamlComponent arch.Component,
	yamlDocument arch.Document,
	mayDependOn []speca.Referable[string],
) error {
	projectImports, err := m.allowedProjectImportsAssembler.assemble(yamlDocument, unwrap(mayDependOn))
	if err != nil {
		return fmt.Errorf("failed to assemble component project imports: %w", err)
	}

	cmp.AllowedProjectImports = wrap(yamlComponent.Reference(), projectImports)
	return nil
}

func (m *componentsAssembler) enrichWithVendorGlobs(
	cmp *speca.Component,
	yamlDocument arch.Document,
	canUse []speca.Referable[string],
) error {
	vendorGlobs, err := m.allowedVendorImportsAssembler.assemble(yamlDocument, unwrap(canUse))
	if err != nil {
		return fmt.Errorf("failed to assemble component vendor imports: %w", err)
	}

	cmp.AllowedVendorGlobs = vendorGlobs
	return nil
}
