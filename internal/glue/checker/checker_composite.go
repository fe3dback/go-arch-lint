package checker

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type CompositeChecker struct {
	checkers []checker
}

func NewCompositeChecker(checkers ...checker) *CompositeChecker {
	return &CompositeChecker{checkers: checkers}
}

func (c *CompositeChecker) Check(spec speca.Spec) (models.CheckResult, error) {
	overallResults := models.CheckResult{}

	for _, checker := range c.checkers {
		results, err := checker.Check(spec)
		if err != nil {
			return models.CheckResult{}, fmt.Errorf("checker failed '%T': %w", checker, err)
		}

		overallResults.Append(results)
	}

	return overallResults, nil
}
