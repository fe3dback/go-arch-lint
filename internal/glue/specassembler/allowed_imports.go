package specassembler

import (
	"fmt"

	yaml "github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
	"github.com/fe3dback/go-arch-lint/internal/models"
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
	spec *yaml.YamlSpec,
	componentNames []models.ComponentName,
	vendorNames []models.VendorName,
) ([]*models.ResolvedPath, error) {
	list := make([]*models.ResolvedPath, 0)

	allowedComponents := make([]models.ComponentName, 0)
	allowedComponents = append(allowedComponents, componentNames...)
	allowedComponents = append(allowedComponents, spec.CommonComponents...)

	allowedVendors := make([]models.VendorName, 0)
	allowedVendors = append(allowedVendors, vendorNames...)
	allowedVendors = append(allowedVendors, spec.CommonVendors...)

	for _, name := range allowedComponents {
		maskPath := spec.Components[name].LocalPath

		resolved, err := aia.resolver.resolveLocalPath(maskPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve mask '%s'", maskPath)
		}

		for _, resolvedPath := range resolved {
			list = append(list, resolvedPath)
		}
	}

	for _, name := range allowedVendors {
		importPath := spec.Vendors[name].ImportPath
		localPath := fmt.Sprintf("vendor/%s", importPath)
		absPath := fmt.Sprintf("%s/%s", aia.rootDirectory, localPath)

		list = append(list, &models.ResolvedPath{
			ImportPath: importPath,
			LocalPath:  localPath,
			AbsPath:    absPath,
		})
	}

	return list, nil
}
