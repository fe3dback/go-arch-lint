package xpath

import (
	"fmt"
	pathUtils "path"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type MatcherAbsolute struct {
	matcherRelative typeMatcher
}

func NewMatcherAbsolute(
	matcherRelative typeMatcher,
) *MatcherAbsolute {
	return &MatcherAbsolute{
		matcherRelative: matcherRelative,
	}
}

func (m *MatcherAbsolute) match(ctx *queryContext, query models.FileQuery) ([]models.FileDescriptor, error) {
	path := query.Path.(models.PathAbsolute) // guaranteed by root composite
	path = models.PathAbsolute(pathUtils.Join(string(query.WorkingDirectory), string(path)))

	// take relative path
	relPathStr, err := filepath.Rel(string(ctx.projectDirectory), string(path))
	if err != nil {
		return nil, fmt.Errorf("failed get relative path from '%s': %w", path, err)
	}

	// pass to next matcher
	relPath := models.PathRelative(relPathStr)
	relQuery := query
	relQuery.Path = relPath

	return m.matcherRelative.match(ctx, relQuery)
}
