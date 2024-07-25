package yaml

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func transformV4(tCtx TransformContext, doc ModelV4) models.Config {
	return models.Config{
		Version:          models.NewRef(models.ConfigVersion(doc.Version), tCtx.createReference("$.version")),
		WorkingDirectory: models.NewRef(models.PathRelative(doc.WorkingDirectory), tCtx.createReference("$.workingDirectory")),
		Settings: models.ConfigSettings{
			DeepScan: models.NewRef(true, models.NewInvalidReference()),
			Imports: models.ConfigSettingsImports{
				StrictMode:            models.NewRef(doc.Settings.Imports.StrictMode, tCtx.createReference("$.settings.imports.strictMode")),
				AllowAnyVendorImports: models.NewRef(doc.Settings.Imports.AllowAnyVendorImports, tCtx.createReference("$.settings.imports.allowAnyVendorImports")),
			},
			Tags: transformV4SettingsTags(tCtx, doc.Settings.Tags),
		},
		Exclude: models.ConfigExclude{
			RelativeDirectories: sliceValuesAutoRef(tCtx, doc.Exclude.RelativeDirectories, "$.exclude.directories", func(dir string) models.PathRelative {
				return models.PathRelative(dir)
			}),
			RelativeFiles: sliceValuesAutoRef(tCtx, doc.Exclude.RelativeFiles, "$.exclude.files", func(file string) models.PathRelativeRegExp {
				return models.PathRelativeRegExp(file)
			}),
		},
		Components: models.ConfigComponents{
			Map: mapValuesAutoRef(tCtx, doc.Components, "$.components",
				func(cCtx TransformContext, name string, component ModelV4Component, refBasePath string) (models.ComponentName, models.ConfigComponent) {
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
				func(context TransformContext, name string, vendor ModelV4Vendor, refBasePath string) (models.VendorName, models.ConfigVendor) {
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
			Map: mapValuesAutoRef(tCtx, doc.Dependencies, "$.dependencies",
				func(context TransformContext, cmpName string, deps ModelV4ComponentDependencies, refBasePath string) (models.ComponentName, models.ConfigComponentDependencies) {
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
						CanContainTags: sliceValuesAutoRef(tCtx, deps.CanContainTags, fmt.Sprintf("%s.canContainTags", refBasePath),
							func(tag string) models.StructTag {
								return models.StructTag(tag)
							}),
					}
				}),
		},
	}
}

func transformV4SettingsTags(tCtx TransformContext, tags ModelV4SettingsTags) models.ConfigSettingsTags {
	refBasePath := "$.settings.structTags.allowed"
	ref := tCtx.createReference(refBasePath)

	if len(tags.Allowed) == 0 {
		return models.ConfigSettingsTags{
			Allowed: models.NewRef(models.ConfigSettingsTagsEnumAll, ref),
		}
	}

	if tags.Allowed[0] == "true" {
		return models.ConfigSettingsTags{
			Allowed: models.NewRef(models.ConfigSettingsTagsEnumAll, ref),
		}
	}

	if tags.Allowed[0] == "false" {
		return models.ConfigSettingsTags{
			Allowed: models.NewRef(models.ConfigSettingsTagsEnumNone, ref),
		}
	}

	return models.ConfigSettingsTags{
		Allowed: models.NewRef(models.ConfigSettingsTagsEnumList, ref),
		AllowedList: sliceValuesAutoRef(tCtx, tags.Allowed, refBasePath,
			func(value string) models.StructTag {
				return models.StructTag(value)
			}),
	}
}
