package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
)

type allowedImportsAssembler struct {
	rootDirectory string
	resolver      *resolver
}

func newAllowedImportsAssembler(
	rootDirectory string,
	resolver *resolver,
) *allowedImportsAssembler {
	return &allowedImportsAssembler{
		rootDirectory: rootDirectory,
		resolver:      resolver,
	}
}

func (aia *allowedImportsAssembler) assemble(
	yamlDocument arch.Document,
	componentNames []string,
	vendorNames []string,
) ([]models.ResolvedPath, error) {
	list := make([]models.ResolvedPath, 0)

	allowedComponents := make([]string, 0)
	allowedComponents = append(allowedComponents, componentNames...)
	for _, componentName := range yamlDocument.CommonComponents().List() {
		allowedComponents = append(allowedComponents, componentName.Value())
	}

	allowedVendors := make([]string, 0)
	allowedVendors = append(allowedVendors, vendorNames...)
	for _, vendorName := range yamlDocument.CommonVendors().List() {
		allowedVendors = append(allowedVendors, vendorName.Value())
	}

	for _, name := range allowedComponents {
		yamlComponent, ok := yamlDocument.Components().Map()[name]
		if !ok {
			continue
		}

		for _, componentIn := range yamlComponent.RelativePaths() {
			relativeGlobPath := componentIn.Value()

			resolved, err := aia.resolver.resolveLocalPath(relativeGlobPath)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve component path '%s'", relativeGlobPath)
			}

			list = append(list, resolved...)
		}
	}

	for _, name := range allowedVendors {
		vendor, ok := yamlDocument.Vendors().Map()[name]
		if !ok {
			continue
		}

		for _, vendorIn := range vendor.ImportPaths() {
			relativeGlobPath := vendorIn.Value()
			localPath := fmt.Sprintf("vendor/%s", relativeGlobPath)

			resolved, err := aia.resolver.resolveVendorPath(localPath)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve vendor path '%s'", relativeGlobPath)
			}

			list = append(list, resolved...)
		}
	}

	return list, nil
}
