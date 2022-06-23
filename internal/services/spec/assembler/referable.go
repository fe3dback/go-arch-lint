package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

func wrap[T any](ref models.Reference, list []T) []speca.Referable[T] {
	res := make([]speca.Referable[T], len(list))

	for ind, path := range list {
		res[ind] = speca.NewReferable(path, ref)
	}

	return res
}

func unwrap[T any](refList []speca.Referable[T]) []T {
	res := make([]T, len(refList))

	for ind, r := range refList {
		res[ind] = r.Value()
	}

	return res
}
