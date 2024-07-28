package yaml

type (
	ModelV3 struct {
		Version            int                                     `json:"version"`
		WorkingDirectory   string                                  `json:"workdir"`
		Allow              ModelV3Allow                            `json:"allow"`
		ExcludeDirectories []string                                `json:"exclude"`
		ExcludeFiles       []string                                `json:"excludeFiles"`
		Components         map[string]ModelV3Component             `json:"components"`
		Vendors            map[string]ModelV3Vendor                `json:"vendors"`
		CommonComponents   []string                                `json:"commonComponents"`
		CommonVendors      []string                                `json:"commonVendors"`
		Dependencies       map[string]ModelV3ComponentDependencies `json:"deps"`
	}

	ModelV3Allow struct {
		DepOnAnyVendor bool `json:"depOnAnyVendor"`
		DeepScan       bool `json:"deepScan"`
	}

	ModelV3Component struct {
		In stringList `json:"in"`
	}

	ModelV3Vendor struct {
		In stringList `json:"in"`
	}

	ModelV3ComponentDependencies struct {
		MayDependOn    []string       `json:"mayDependOn"`
		CanUse         []string       `json:"canUse"`
		AnyVendorDeps  bool           `json:"anyVendorDeps"`
		AnyProjectDeps bool           `json:"anyProjectDeps"`
		DeepScan       optional[bool] `json:"deepScan"`
	}
)
