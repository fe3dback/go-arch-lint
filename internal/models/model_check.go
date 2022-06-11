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
		ComponentName      string    `json:"ComponentName"`
		FileRelativePath   string    `json:"FileRelativePath"`
		FileAbsolutePath   string    `json:"FileAbsolutePath"`
		ResolvedImportName string    `json:"ResolvedImportName"`
		Reference          Reference `json:"-"`
		SourceCodePreview  []byte    `json:"-"`
	}

	CheckArchWarningMatch struct {
		FileRelativePath  string    `json:"FileRelativePath"`
		FileAbsolutePath  string    `json:"FileAbsolutePath"`
		Reference         Reference `json:"-"`
		SourceCodePreview []byte    `json:"-"`
	}

	// Example layout:
	// ──────────────────────────────────────────────────────────────────────────
	// Component 'operations', method 'NewOperation' defined in:
	// - internal/glue/code/line_count.go:54
	// Should not depend on 'repository' as 'micro.ViewRepository', injected:
	// - c.provideMicroViewRepository() in:
	//   - internal/app/internal/container/cmd_mapping.go:15
	//
	// 14: operations.NewOperation(
	// 15:   c.provideMicroViewRepository() : micro.ViewRepository
	//       └───────────────┐
	//                       ↓
	// 54: func NewOperation(repo viewRepo) *Operation {
	// 56:   return &Operation{repo: repo}
	// ──────────────────────────────────────────────────────────────────────────

	CheckDeepscanWarning struct {
		Gate       CheckDeepscanWarningGate
		Dependency CheckDeepscanWarningDependency
		LineArt    CheckDeepscanWarningArt
	}

	CheckDeepscanWarningGate struct {
		ComponentName     string    // operations
		MethodName        string    // NewOperation
		Definition        Reference // internal/glue/code/line_count.go:54
		SourceCodePreview []byte    `json:"-"`
	}

	CheckDeepscanWarningDependency struct {
		ComponentName     string    // repository
		Name              string    // micro.ViewRepository
		InjectionAST      string    // c.provideMicroViewRepository()
		Injection         Reference // internal/app/internal/container/cmd_mapping.go:15
		SourceCodePreview []byte    `json:"-"`
	}

	// CheckDeepscanWarningArt used for drawing this:
	//     outpos
	//       └───────────────┐
	//  toRight:Y            ↓
	//  length:15          inpos
	CheckDeepscanWarningArt struct {
		ToRight    bool
		OutPos     int
		InPos      int
		LineLength int
	}

	CheckResult struct {
		DependencyWarnings []CheckArchWarningDependency
		MatchWarnings      []CheckArchWarningMatch
		DeepscanWarnings   []CheckDeepscanWarning
	}
)

func (cr *CheckResult) Append(another CheckResult) {
	cr.DependencyWarnings = append(cr.DependencyWarnings, another.DependencyWarnings...)
	cr.MatchWarnings = append(cr.MatchWarnings, another.MatchWarnings...)
	cr.DeepscanWarnings = append(cr.DeepscanWarnings, another.DeepscanWarnings...)
}
