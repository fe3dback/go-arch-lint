package yaml

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func transformV3(tCtx TransformContext, doc ModelV3) models.Config {
	return models.Config{
		Version:          models.NewRef(models.ConfigVersion(doc.Version), tCtx.createReference("$.version")),
		WorkingDirectory: models.NewRef(models.PathRelative(doc.WorkingDirectory), tCtx.createReference("$.workdir")),
		Settings: models.ConfigSettings{
			DeepScan: models.NewRef(doc.Allow.DeepScan, tCtx.createReference("$.allow.deepScan")),
			Imports: models.ConfigSettingsImports{
				StrictMode:            models.NewRef(false, models.NewInvalidReference()),
				AllowAnyVendorImports: models.NewRef(doc.Allow.DepOnAnyVendor, tCtx.createReference("$.allow.depOnAnyVendor")),
			},
			Tags: models.ConfigSettingsTags{
				Allowed: models.NewRef(models.ConfigSettingsTagsEnumAll, models.NewInvalidReference()),
			},
		},
		Exclude: models.ConfigExclude{
			RelativeDirectories: sliceValuesAutoRef(tCtx, doc.ExcludeDirectories, "$.exclude", func(dir string) models.PathRelative {
				return models.PathRelative(dir)
			}),
			RelativeFiles: sliceValuesAutoRef(tCtx, doc.ExcludeFiles, "$.excludeFiles", func(file string) models.PathRelativeRegExp {
				return models.PathRelativeRegExp(file)
			}),
		},
		Components: models.ConfigComponents{
			Map: mapValuesAutoRef(tCtx, doc.Components, "$.components",
				func(tCtx TransformContext, name string, component ModelV3Component, refBasePath string) (models.ComponentName, models.ConfigComponent) {
					return models.ComponentName(name), models.ConfigComponent{
						In: sliceValuesAutoRef(tCtx, component.In, fmt.Sprintf("%s.in", refBasePath),
							func(value string) models.PathRelativeGlob {
								return models.PathRelativeGlob(value)
							}),
					}
				}),
		},
		Vendors: models.ConfigVendors{
			Map: mapValuesAutoRef(tCtx, doc.Vendors, "$.vendors",
				func(tCtx TransformContext, name string, vendor ModelV3Vendor, refBasePath string) (models.VendorName, models.ConfigVendor) {
					return models.VendorName(name), models.ConfigVendor{
						In: sliceValuesAutoRef(tCtx, vendor.In, fmt.Sprintf("%s.in", refBasePath),
							func(value string) models.PathImportGlob {
								return models.PathImportGlob(value)
							}),
					}
				}),
		},
		CommonComponents: sliceValuesAutoRef(tCtx, doc.CommonComponents, "$.commonComponents", func(v string) models.ComponentName {
			return models.ComponentName(v)
		}),
		CommonVendors: sliceValuesAutoRef(tCtx, doc.CommonVendors, "$.commonVendors", func(v string) models.VendorName {
			return models.VendorName(v)
		}),
		Dependencies: models.ConfigDependencies{
			Map: mapValuesAutoRef(tCtx, doc.Dependencies, "$.deps",
				func(tCtx TransformContext, cmpName string, deps ModelV3ComponentDependencies, refBasePath string) (models.ComponentName, models.ConfigComponentDependencies) {
					return models.ComponentName(cmpName), models.ConfigComponentDependencies{
						MayDependOn: sliceValuesAutoRef(tCtx, deps.MayDependOn, fmt.Sprintf("%s.mayDependOn", refBasePath),
							func(anotherCmpName string) models.ComponentName {
								return models.ComponentName(anotherCmpName)
							}),
						CanUse: sliceValuesAutoRef(tCtx, deps.CanUse, fmt.Sprintf("%s.canUse", refBasePath),
							func(vendorName string) models.VendorName {
								return models.VendorName(vendorName)
							}),
						AnyVendorDeps:  models.NewRef(deps.AnyVendorDeps, tCtx.createReference(fmt.Sprintf("%s.anyVendorDeps", refBasePath))),
						AnyProjectDeps: models.NewRef(deps.AnyProjectDeps, tCtx.createReference(fmt.Sprintf("%s.anyProjectDeps", refBasePath))),
						CanContainTags: []models.Ref[models.StructTag]{},
					}
				}),
		},
	}
}
