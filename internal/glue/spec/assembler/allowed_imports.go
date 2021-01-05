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

		maskPath := yamlComponent.LocalPath().Value()

		resolved, err := aia.resolver.resolveLocalPath(maskPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve mask '%s'", maskPath)
		}

		list = append(list, resolved...)
	}

	for _, name := range allowedVendors {
		vendor, ok := yamlDocument.Vendors().Map()[name]
		if !ok {
			continue
		}

		for _, vendorIn := range vendor.ImportPaths() {
			relativeGlobPath := vendorIn.Value()
			localPath := fmt.Sprintf("vendor/%s", relativeGlobPath)

			resolvedPathList, err := aia.resolver.resolveVendorPath(localPath)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve vendor path '%s'", relativeGlobPath)
			}

			list = append(list, resolvedPathList...)
		}
	}

	return list, nil
}
