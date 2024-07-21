package yaml

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

func transformV4(tCtx TransformContext, doc ModelV4) models.Config {
	return models.Config{
		Version:          models.NewRef(models.ConfigVersion(doc.Version), tCtx.createReference("$.version")),
		WorkingDirectory: models.NewRef(models.PathRelative(doc.WorkingDirectory), tCtx.createReference("$.workingDirectory")),
		Settings:         models.ConfigSettings{},
		Exclude:          models.ConfigExclude{},
		Components:       models.ConfigComponents{},
		Vendors:          models.ConfigVendors{},
		CommonComponents: nil,
		CommonVendors:    nil,
		Dependencies:     models.ConfigDependencies{},
	}
}
