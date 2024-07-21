package container

import "github.com/fe3dback/go-arch-lint/v4/internal/operations/mapping"

func (c *Container) operationMapping() *mapping.Operation {
	return once(func() *mapping.Operation {
		return mapping.NewOperation(
			c.serviceConfigReader(),
		)
	})
}
