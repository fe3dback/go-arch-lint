package checker

import (
	"context"
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

func (c *CompositeChecker) Check(ctx context.Context, spec speca.Spec) (models.CheckResult, error) {
	overallResults := models.CheckResult{}

	for ind, checker := range c.checkers {
		results, err := checker.Check(ctx, spec)
		if err != nil {
			return models.CheckResult{}, fmt.Errorf("checker failed '%T': %w", checker, err)
		}

		overallResults.Append(results)

		if results.HasNotices() && ind < len(c.checkers)-1 {
			fmt.Printf("skipped other checks, found early lint notices..\n\n")
			break
		}
	}

	return overallResults, nil
}
