package models

type (
	Check struct {
		DocumentNotices        []CheckNotice                `json:"ExecutionWarnings"`
		ArchHasWarnings        bool                         `json:"ArchHasWarnings"`
		ArchWarningsDependency []CheckArchWarningDependency `json:"ArchWarningsDeps"`
		ArchWarningsMatch      []CheckArchWarningMatch      `json:"ArchWarningsNotMatched"`
		ModuleName             string                       `json:"ModuleName"`
	}

	CheckNotice struct {
		Text              string `json:"Text"`
		File              string `json:"File"`
		Line              int    `json:"Line"`
		Offset            int    `json:"Offset"`
		SourceCodePreview []byte `json:"-"`
	}

	CheckArchWarningDependency struct {
		ComponentName      string `json:"ComponentName"`
		FileRelativePath   string `json:"FileRelativePath"`
		FileAbsolutePath   string `json:"FileAbsolutePath"`
		ResolvedImportName string `json:"ResolvedImportName"`
	}

	CheckArchWarningMatch struct {
		FileRelativePath string `json:"FileRelativePath"`
		FileAbsolutePath string `json:"FileAbsolutePath"`
	}
)
