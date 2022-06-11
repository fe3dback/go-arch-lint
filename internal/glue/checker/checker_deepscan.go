package checker

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type DeepScan struct {
}

func NewDeepScan() *DeepScan {
	return &DeepScan{}
}

func (c *DeepScan) Check(spec speca.Spec) (models.CheckResult, error) {
	overallResults := models.CheckResult{}

	for _, component := range spec.Components {
		results, err := c.checkComponent(component)
		if err != nil {
			return models.CheckResult{}, fmt.Errorf("component '%s' check failed: %w",
				component.Name.Value(),
				err,
			)
		}

		overallResults.Append(results)
	}

	return overallResults, nil
}

func (c *DeepScan) checkComponent(cmp speca.Component) (models.CheckResult, error) {
	// todo check

	return models.CheckResult{}, nil
}
