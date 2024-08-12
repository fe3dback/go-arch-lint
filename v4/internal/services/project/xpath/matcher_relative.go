package xpath

import (
	pathUtils "path"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type MatcherRelative struct {
}

func NewMatcherRelative() *MatcherRelative {
	return &MatcherRelative{}
}

func (m *MatcherRelative) match(ctx *queryContext, query models.FileQuery) ([]models.FileDescriptor, error) {
	path := query.Path.(models.PathRelative) // guaranteed by root composite
	path = models.PathRelative(pathUtils.Join(string(query.WorkingDirectory), string(path)))

	descriptors := make([]models.FileDescriptor, 0, 2)

	if query.Type == models.FileMatchQueryTypeAll || query.Type == models.FileMatchQueryTypeOnlyDirectories {
		if dir, found := ctx.index.directoryAt(path); found {
			descriptors = append(descriptors, dir)
		}
	}

	if query.Type == models.FileMatchQueryTypeAll || query.Type == models.FileMatchQueryTypeOnlyFiles {
		if file, found := ctx.index.fileAt(path); found {
			descriptors = append(descriptors, file)
		}
	}

	return descriptors, nil
}
