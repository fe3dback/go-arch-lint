package specassembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

func wrapPaths(ref speca.Reference, list []models.ResolvedPath) []speca.ReferableResolvedPath {
	res := make([]speca.ReferableResolvedPath, len(list))

	for ind, path := range list {
		res[ind] = speca.NewReferableResolvedPath(
			path,
			ref,
		)
	}

	return res
}

func unwrapStrings(rs []speca.ReferableString) []string {
	res := make([]string, len(rs))

	for ind, r := range rs {
		res[ind] = r.Value()
	}

	return res
}
