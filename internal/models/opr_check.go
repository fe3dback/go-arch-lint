package models

import "github.com/fe3dback/go-arch-lint/internal/models/common"

type (
	CmdCheckIn struct {
		ProjectPath string
		ArchFile    string
		MaxWarnings int
	}

	CmdCheckOut struct {
		DocumentNotices        []CheckNotice                `json:"ExecutionWarnings"`
		ArchHasWarnings        bool                         `json:"ArchHasWarnings"`
		ArchWarningsDependency []CheckArchWarningDependency `json:"ArchWarningsDeps"`
		ArchWarningsMatch      []CheckArchWarningMatch      `json:"ArchWarningsNotMatched"`
		ArchWarningsDeepScan   []CheckArchWarningDeepscan   `json:"ArchWarningsDeepScan"`
		OmittedCount           int                          `json:"OmittedCount"`
		ModuleName             string                       `json:"ModuleName"`
	}

	CheckNotice struct {
		Text              string `json:"Text"`
		File              string `json:"File"`
		Line              int    `json:"Line"`
		Column            int    `json:"Offset"`
		SourceCodePreview []byte `json:"-"`
	}

	CheckArchWarningDependency struct {
		ComponentName      string           `json:"ComponentName"`
		FileRelativePath   string           `json:"FileRelativePath"`
		FileAbsolutePath   string           `json:"FileAbsolutePath"`
		ResolvedImportName string           `json:"ResolvedImportName"`
		Reference          common.Reference `json:"Reference"`
	}

	CheckArchWarningMatch struct {
		FileRelativePath string           `json:"FileRelativePath"`
		FileAbsolutePath string           `json:"FileAbsolutePath"`
		Reference        common.Reference `json:"-"`
	}

	CheckArchWarningDeepscan struct {
		Gate       DeepscanWarningGate       `json:"Gate"`
		Dependency DeepscanWarningDependency `json:"Dependency"`
		Target     DeepscanWarningTarget     `json:"Target"`
	}

	DeepscanWarningGate struct {
		ComponentName string           `json:"ComponentName"` // operations
		MethodName    string           `json:"MethodName"`    // NewOperation
		Definition    common.Reference `json:"Definition"`    // internal/glue/code/line_count.go:54
		RelativePath  string           `json:"-"`             // internal/glue/code/line_count.go:54
	}

	DeepscanWarningDependency struct {
		ComponentName     string           `json:"ComponentName"` // repository
		Name              string           `json:"Name"`          // micro.ViewRepository
		InjectionAST      string           `json:"InjectionAST"`  // c.provideMicroViewRepository()
		Injection         common.Reference `json:"Injection"`     // internal/app/internal/container/cmd_mapping.go:15
		InjectionPath     string           `json:"-"`             // internal/app/internal/container/cmd_mapping.go:15
		SourceCodePreview []byte           `json:"-"`
	}

	DeepscanWarningTarget struct {
		Definition   common.Reference `json:"Definition"`
		RelativePath string           `json:"-"` // internal/app/internal/container/cmd_mapping.go:15
	}

	CheckResult struct {
		DependencyWarnings []CheckArchWarningDependency
		MatchWarnings      []CheckArchWarningMatch
		DeepscanWarnings   []CheckArchWarningDeepscan
	}
)

func (cr *CheckResult) Append(another CheckResult) {
	cr.DependencyWarnings = append(cr.DependencyWarnings, another.DependencyWarnings...)
	cr.MatchWarnings = append(cr.MatchWarnings, another.MatchWarnings...)
	cr.DeepscanWarnings = append(cr.DeepscanWarnings, another.DeepscanWarnings...)
}

func (cr *CheckResult) HasNotices() bool {
	if len(cr.DependencyWarnings) > 0 {
		return true
	}
	if len(cr.MatchWarnings) > 0 {
		return true
	}
	if len(cr.DeepscanWarnings) > 0 {
		return true
	}

	return false
}
