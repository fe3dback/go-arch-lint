package yaml

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func sliceValuesAutoRef[V any, T any](tCtx TransformContext, in []V, baseRef string, fn func(V) T) models.RefSlice[T] {
	result := make(models.RefSlice[T], 0, len(in))

	if len(in) == 1 {
		ref := tCtx.createReference(fmt.Sprintf("%s", baseRef))
		result = append(result, models.NewRef(fn(in[0]), ref))

		return result
	}

	for ind, value := range in {
		ref := tCtx.createReference(fmt.Sprintf("%s[%d]", baseRef, ind))
		result = append(result, models.NewRef(fn(value), ref))
	}

	return result
}

func mapValuesAutoRef[K comparable, V any, TK comparable, TV any](
	tCtx TransformContext, in map[K]V, baseRef string, fn func(TransformContext, K, V, string) (TK, TV)) models.RefMap[TK, TV] {

	result := models.NewRefMap[TK, TV](len(in))

	for key, value := range in {
		keyPath := fmt.Sprintf("%s.%v", baseRef, key)

		tk, tv := fn(tCtx, key, value, keyPath)
		ref := tCtx.createReference(keyPath)

		result.Set(tk, tv, ref)
	}

	return result
}
