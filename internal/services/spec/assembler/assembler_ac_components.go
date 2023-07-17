package assembler

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/services/spec"
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

func (m *componentsAssembler) assemble(spec *arch.Spec, document spec.Document) error {
	for yamlName, yamlComponent := range document.Components() {
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
	yamlComponent common.Referable[spec.Component],
	yamlDocument spec.Document,
) (arch.Component, error) {
	depMeta, hasDeps := yamlDocument.Dependencies()[yamlName]

	mayDependOn := make([]common.Referable[string], 0)
	canUse := make([]common.Referable[string], 0)
	deepScan := yamlDocument.Options().DeepScan()

	if hasDeps {
		mayDependOn = append(mayDependOn, depMeta.Value.MayDependOn()...)
		canUse = append(canUse, depMeta.Value.CanUse()...)
		deepScan = depMeta.Value.DeepScan()
	}

	cmp := arch.Component{
		Name:        common.NewReferable(yamlName, yamlComponent.Reference),
		MayDependOn: mayDependOn,
		CanUse:      canUse,
		DeepScan:    deepScan,
	}

	type enricher func() error
	enrichers := []enricher{
		func() error { return m.enrichWithFlags(&cmp, yamlComponent, hasDeps, depMeta.Value) },
		func() error { return m.enrichWithResolvedPaths(&cmp, yamlDocument, yamlName, yamlComponent) },
		func() error { return m.enrichWithProjectImports(&cmp, yamlComponent, yamlDocument, mayDependOn) },
		func() error { return m.enrichWithVendorGlobs(&cmp, yamlDocument, canUse) },
	}

	for _, enrich := range enrichers {
		err := enrich()
		if err != nil {
			return arch.Component{}, fmt.Errorf("failed assemble component '%s', enrich '%T' err: %w",
				yamlName,
				enrich,
				err,
			)
		}
	}

	return cmp, nil
}

func (m *componentsAssembler) enrichWithFlags(
	cmp *arch.Component,
	yamlComponent common.Referable[spec.Component],
	hasDeps bool,
	depMeta spec.DependencyRule,
) error {
	if hasDeps {
		cmp.SpecialFlags = arch.SpecialFlags{
			AllowAllProjectDeps: depMeta.AnyProjectDeps(),
			AllowAllVendorDeps:  depMeta.AnyVendorDeps(),
		}
		return nil
	}

	cmp.SpecialFlags = arch.SpecialFlags{
		AllowAllProjectDeps: common.NewReferable(false, yamlComponent.Reference),
		AllowAllVendorDeps:  common.NewReferable(false, yamlComponent.Reference),
	}

	return nil
}

func (m *componentsAssembler) enrichWithResolvedPaths(
	cmp *arch.Component,
	yamlDocument spec.Document,
	yamlName string,
	yamlComponent common.Referable[spec.Component],
) error {
	resolvedPaths := make([]common.Referable[models.ResolvedPath], 0)

	for _, componentIn := range yamlComponent.Value.RelativePaths() {
		tmpResolvedPath, err := m.resolver.resolveLocalGlobPath(
			path.Clean(fmt.Sprintf("%s/%s",
				yamlDocument.WorkingDirectory().Value,
				string(componentIn),
			)),
		)
		if err != nil {
			return fmt.Errorf("failed to assemble component '%s' path '%s': %w",
				yamlName,
				componentIn,
				err,
			)
		}

		wrappedPaths := wrap(yamlComponent.Reference, tmpResolvedPath)
		resolvedPaths = append(resolvedPaths, wrappedPaths...)
	}

	cmp.ResolvedPaths = resolvedPaths
	return nil
}

func (m *componentsAssembler) enrichWithProjectImports(
	cmp *arch.Component,
	yamlComponent common.Referable[spec.Component],
	yamlDocument spec.Document,
	mayDependOn []common.Referable[string],
) error {
	projectImports, err := m.allowedProjectImportsAssembler.assemble(yamlDocument, unwrap(mayDependOn))
	if err != nil {
		return fmt.Errorf("failed to assemble component project imports: %w", err)
	}

	cmp.AllowedProjectImports = wrap(yamlComponent.Reference, projectImports)
	return nil
}

func (m *componentsAssembler) enrichWithVendorGlobs(
	cmp *arch.Component,
	yamlDocument spec.Document,
	canUse []common.Referable[string],
) error {
	vendorGlobs, err := m.allowedVendorImportsAssembler.assemble(yamlDocument, unwrap(canUse))
	if err != nil {
		return fmt.Errorf("failed to assemble component vendor imports: %w", err)
	}

	cmp.AllowedVendorGlobs = vendorGlobs
	return nil
}
