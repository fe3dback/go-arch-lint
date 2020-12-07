package checker

import (
	"sort"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type (
	results models.CheckResult
)

func newResults() results {
	return results{
		DependencyWarnings: []models.CheckArchWarningDependency{},
		MatchWarnings:      []models.CheckArchWarningMatch{},
	}
}

func (res *results) addNotMatchedWarning(warn models.CheckArchWarningMatch) {
	res.MatchWarnings = append(res.MatchWarnings, warn)
}

func (res *results) addDependencyWarning(warn models.CheckArchWarningDependency) {
	res.DependencyWarnings = append(res.DependencyWarnings, warn)
}

func (res *results) assembleSortedResults() models.CheckResult {
	sort.Slice(res.DependencyWarnings, func(i, j int) bool {
		return res.DependencyWarnings[i].FileRelativePath < res.DependencyWarnings[j].FileRelativePath
	})

	sort.Slice(res.MatchWarnings, func(i, j int) bool {
		return res.MatchWarnings[i].FileRelativePath < res.MatchWarnings[j].FileRelativePath
	})

	return models.CheckResult{
		DependencyWarnings: res.DependencyWarnings,
		MatchWarnings:      res.MatchWarnings,
	}
}
