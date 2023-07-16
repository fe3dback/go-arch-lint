package assembler

import (
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

func wrap[T any](ref common.Reference, list []T) []common.Referable[T] {
	res := make([]common.Referable[T], len(list))

	for ind, path := range list {
		res[ind] = common.NewReferable(path, ref)
	}

	return res
}

func unwrap[T any](refList []common.Referable[T]) []T {
	res := make([]T, len(refList))

	for ind, r := range refList {
		res[ind] = r.Value
	}

	return res
}
