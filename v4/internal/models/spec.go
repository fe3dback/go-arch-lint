package models

type (
	ComponentName string

	VendorName string

	Spec struct {
		ProjectDirectory PathAbsolute
		WorkingDirectory PathRelative
		ModuleName       string
		Components       []SpecComponent
	}

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
		OwnedFiles          []PathRelative
		OwnedPackages       []PathRelative
	}
)
