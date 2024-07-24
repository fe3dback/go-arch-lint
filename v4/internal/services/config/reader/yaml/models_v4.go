package yaml

type (
	ModelV4 struct {
		Version          int                                     `json:"version"`
		WorkingDirectory string                                  `json:"workingDirectory"`
		Settings         ModelV4Settings                         `json:"settings"`
		Exclude          ModelV4Exclude                          `json:"exclude"`
		Components       map[string]ModelV4Component             `json:"components"`
		Vendors          map[string]ModelV4Vendor                `json:"vendors"`
		CommonComponents []string                                `json:"commonComponents"`
		CommonVendors    []string                                `json:"commonVendors"`
		Dependencies     map[string]ModelV4ComponentDependencies `json:"dependencies"`
	}

	ModelV4Settings struct {
		Imports ModelV4SettingsImports `json:"imports"`
		Tags    ModelV4SettingsTags    `json:"structTags"`
	}

	ModelV4SettingsImports struct {
		StrictMode            bool `json:"strictMode"`
		AllowAnyVendorImports bool `json:"allowAnyVendorImports"`
	}

	ModelV4SettingsTags struct {
		Allowed stringList `json:"allowed"`
	}

	ModelV4Exclude struct {
		RelativeDirectories []string `json:"directories"`
		RelativeFiles       []string `json:"files"`
	}

	ModelV4Component struct {
		In stringList `json:"in"`
	}

	ModelV4Vendor struct {
		In stringList `json:"in"`
	}

	ModelV4ComponentDependencies struct {
		MayDependOn    []string `json:"mayDependOn"`
		CanUse         []string `json:"canUse"`
		AnyVendorDeps  bool     `json:"anyVendorDeps"`
		AnyProjectDeps bool     `json:"anyProjectDeps"`
		CanContainTags []string `json:"canContainTags"`
	}
)
