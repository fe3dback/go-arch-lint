package resolver

import (
	"regexp"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

func refPathToList(list []speca.Referable[models.ResolvedPath]) []models.ResolvedPath {
	result := make([]models.ResolvedPath, 0)

	for _, path := range list {
		result = append(result, path.Value())
	}

	return result
}

func refRegExpToList(list []speca.Referable[*regexp.Regexp]) []*regexp.Regexp {
	result := make([]*regexp.Regexp, 0)

	for _, path := range list {
		result = append(result, path.Value())
	}

	return result
}
