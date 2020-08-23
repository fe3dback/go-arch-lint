package checker

type (
	CheckResult struct {
		deps  []WarningDep
		match []WarningNotMatched
	}

	WarningDep struct {
		ComponentName      string
		FileRelativePath   string
		FileAbsolutePath   string
		ResolvedImportName string
	}

	WarningNotMatched struct {
		FileRelativePath string
		FileAbsolutePath string
	}
)

func newCheckResult() *CheckResult {
	return &CheckResult{
		deps:  make([]WarningDep, 0),
		match: make([]WarningNotMatched, 0),
	}
}

func (res *CheckResult) IsOk() bool {
	if len(res.deps) != 0 {
		return false
	}

	if len(res.match) != 0 {
		return false
	}

	return true
}

func (res *CheckResult) TotalCount() int {
	return len(res.match) + len(res.deps)
}

func (res *CheckResult) DependencyWarnings() []WarningDep {
	return res.deps
}

func (res *CheckResult) NotMatchedWarnings() []WarningNotMatched {
	return res.match
}

func (res *CheckResult) addNotMatchedWarning(warn WarningNotMatched) {
	res.match = append(res.match, warn)
}

func (res *CheckResult) addDependencyWarning(warn WarningDep) {
	res.deps = append(res.deps, warn)
}
