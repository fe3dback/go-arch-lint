package models

const (
	TagsAllowedEnumAll  ConfigSettingsTagsEnum = "All"
	TagsAllowedEnumNone ConfigSettingsTagsEnum = "None"
	TagsAllowedEnumList ConfigSettingsTagsEnum = "List"
)

type (
	ConfigSettingsTagsEnum string
)

type (
	Config struct {
		SyntaxProblems RefSlice[string]

		Version          Ref[ConfigVersion]
		WorkingDirectory Ref[PathRelative]
		Settings         ConfigSettings
		Exclude          ConfigExclude
		Components       ConfigComponents
		Vendors          ConfigVendors
		CommonComponents []ComponentName
		CommonVendors    []VendorName
		Dependencies     ConfigDependencies
	}

	ConfigSettings struct {
		Imports ConfigSettingsImports
		Tags    ConfigSettingsTags
	}

	ConfigSettingsImports struct {
		StrictMode            Ref[bool]
		AllowAnyVendorImports Ref[bool]
	}

	ConfigSettingsTags struct {
		Allowed     Ref[ConfigSettingsTagsEnum]
		AllowedList RefSlice[StructTag]
	}

	ConfigExclude struct {
		RelativeDirectories RefSlice[PathRelative]
		RelativeFiles       RefSlice[PathRelativeRegExp]
	}

	ConfigComponents struct {
		Map RefMap[ComponentName, ConfigComponent]
	}

	ConfigComponent struct {
		In RefSlice[PathRelativeGlob]
	}

	ConfigVendors struct {
		Map RefMap[VendorName, ConfigVendor]
	}

	ConfigVendor struct {
		In RefSlice[PathImportGlob]
	}

	ConfigDependencies struct {
		Map RefMap[ComponentName, ConfigComponentDependencies]
	}

	ConfigComponentDependencies struct {
		MayDependOn    RefSlice[ComponentName]
		CanUse         RefSlice[VendorName]
		AnyVendorDeps  Ref[bool]
		AnyProjectDeps Ref[bool]
		CanContainTags Ref[bool]
	}
)
