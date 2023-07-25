package decoder

import (
	"context"
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-yaml/ast"
)

type ref[T any] struct {
	defined bool
	ref     common.Referable[T]
}

type stringList []string
type yamlParentFileCtx struct{}

func (r *ref[T]) UnmarshalYAML(ctx context.Context, node ast.Node, decode func(interface{}) error) error {
	filePath := ""
	if ref, ok := ctx.Value(yamlParentFileCtx{}).(string); ok {
		filePath = ref
	}

	r.defined = true
	r.ref.Reference = common.NewReferenceSingleLine(
		filePath,
		node.GetToken().Position.Line,
		node.GetToken().Position.Column,
	)

	return decode(&r.ref.Value)
}

func (s *stringList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []string
	var lastErr error

	if err := unmarshal(&list); err == nil {
		*s = list
		return nil
	} else {
		lastErr = err
	}

	var value string
	if err := unmarshal(&value); err == nil {
		*s = []string{value}
		return nil
	} else {
		lastErr = fmt.Errorf("%v: %w", lastErr, err)
	}

	return fmt.Errorf("failed decode yaml stringsList: %w", lastErr)
}
