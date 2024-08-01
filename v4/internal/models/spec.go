package models

type (
	ComponentName string

	VendorName string

	Spec struct {
		Project          ProjectInfo
		WorkingDirectory Ref[PathRelative]
		Components       []SpecComponent
	}

	// todo: spec checking
	// - pick file
	// - find parent component
	// - can check imports already
	// - can check strictMode already
	// - todo: some cache system (3 files: gate, injector, dependency)??

	SpecComponent struct {
		Name                Ref[ComponentName]
		DefinitionComponent Reference // $.components.<NAME>
		DefinitionDeps      Reference // $.deps.<NAME>
		DeepScan            Ref[bool]
		StrictMode          Ref[bool]
		AllowAllProjectDeps Ref[bool]
		AllowAllVendorDeps  Ref[bool]
		AllowAllTags        Ref[bool]
		AllowedTags         RefSlice[StructTag]
		MayDependOn         RefSlice[ComponentName]
		CanUse              RefSlice[VendorName]
		MatchedFiles        []PathRelative // all files matched by component "in" query
		MatchedPackages     []PathRelative // all packages matched by component "in" query
		OwnedFiles          []PathRelative // unique subset of MatchedFiles, belongs to this component (every file will belong only to single component)
		OwnedPackages       []PathRelative // unique subset of MatchedPackages, belongs to this component (every package will belong only to single component)
	}
)
