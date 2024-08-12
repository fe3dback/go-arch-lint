package xpath

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type queryContext struct {
	projectDirectory models.PathAbsolute
	index            *index
}

func newQueryContext(projectDirectory models.PathAbsolute) queryContext {
	return queryContext{
		projectDirectory: projectDirectory,
		index:            newIndex(),
	}
}
