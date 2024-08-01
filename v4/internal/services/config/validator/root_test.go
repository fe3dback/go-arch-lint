package validator

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

const (
	cmpContainer  = models.ComponentName("container")
	cmpHandler    = models.ComponentName("handler")
	cmpUsecase    = models.ComponentName("usecase")
	cmpRepository = models.ComponentName("repository")
	cmpModels     = models.ComponentName("models")

	vendorDB      = models.VendorName("db")
	vendorTracing = models.VendorName("tracing")
)

type validatorIn struct {
	conf models.Config
}

// shortcut for models.NewInvalidRef
// for better code readability
func nf[T any](value T) models.Ref[T] {
	return models.NewInvalidRef(value)
}

//nolint:funlen
func createValidatorIn(mutators ...func(*validatorIn)) validatorIn {
	in := validatorIn{
		conf: models.Config{
			Version:          nf(models.ConfigVersion(4)),
			WorkingDirectory: nf(models.PathRelative("internal")),
			Settings: models.ConfigSettings{
				DeepScan: nf(false),
				Imports: models.ConfigSettingsImports{
					StrictMode:            nf(false),
					AllowAnyVendorImports: nf(false),
				},
				Tags: models.ConfigSettingsTags{
					Allowed: nf(models.ConfigSettingsTagsEnumAll),
				},
			},
			Components: models.ConfigComponents{
				Map: models.NewRefMapFrom(map[models.ComponentName]models.Ref[models.ConfigComponent]{
					cmpContainer: nf(models.ConfigComponent{
						In: []models.Ref[models.PathRelativeGlob]{
							nf(models.PathRelativeGlob("app/internal/container")),
							nf(models.PathRelativeGlob("plugins/*/container")),
						},
					}),
					cmpHandler: nf(models.ConfigComponent{
						In: []models.Ref[models.PathRelativeGlob]{
							nf(models.PathRelativeGlob("handlers")),
						},
					}),
					cmpUsecase: nf(models.ConfigComponent{
						In: []models.Ref[models.PathRelativeGlob]{
							nf(models.PathRelativeGlob("app/*/business/**")),
						},
					}),
					cmpRepository: nf(models.ConfigComponent{
						In: []models.Ref[models.PathRelativeGlob]{
							nf(models.PathRelativeGlob("app/*/repo")),
						},
					}),
					cmpModels: nf(models.ConfigComponent{
						In: []models.Ref[models.PathRelativeGlob]{
							nf(models.PathRelativeGlob("models/**")),
						},
					}),
				}),
			},
			Vendors: models.ConfigVendors{
				Map: models.NewRefMapFrom(map[models.VendorName]models.Ref[models.ConfigVendor]{
					vendorDB: nf(models.ConfigVendor{
						In: []models.Ref[models.PathImportGlob]{
							nf(models.PathImportGlob("github.com/fe3dback/orm")),
							nf(models.PathImportGlob("github.com/fe3dback/libs/**/db/*")),
							nf(models.PathImportGlob("github.com/fe3dback/transactional")),
						},
					}),
					vendorTracing: nf(models.ConfigVendor{
						In: []models.Ref[models.PathImportGlob]{
							nf(models.PathImportGlob("io.org.example.com/telemetry/*/tracing")),
						},
					}),
				}),
			},
			CommonComponents: []models.Ref[models.ComponentName]{
				nf(cmpModels),
			},
			CommonVendors: []models.Ref[models.VendorName]{
				nf(vendorTracing),
			},
			Dependencies: models.ConfigDependencies{
				Map: models.NewRefMapFrom(map[models.ComponentName]models.Ref[models.ConfigComponentDependencies]{
					cmpContainer: nf(models.ConfigComponentDependencies{
						AnyVendorDeps:  nf(true),
						AnyProjectDeps: nf(true),
					}),
					cmpHandler: nf(models.ConfigComponentDependencies{
						MayDependOn: []models.Ref[models.ComponentName]{
							nf(cmpUsecase),
						},
					}),
					cmpUsecase: nf(models.ConfigComponentDependencies{
						MayDependOn: []models.Ref[models.ComponentName]{
							nf(cmpUsecase),
							nf(cmpRepository),
						},
					}),
					cmpRepository: nf(models.ConfigComponentDependencies{
						CanUse: []models.Ref[models.VendorName]{
							nf(vendorDB),
						},
						CanContainTags: []models.Ref[models.StructTag]{
							nf(models.StructTag("db")),
						},
					}),
				}),
			},
		},
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}
