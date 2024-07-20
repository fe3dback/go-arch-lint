package ptr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fe3dback/go-arch-lint/v4/internal/pkg/ptr"
)

func TestRef(t *testing.T) {
	val1 := 50

	assert.Equal(t, &val1, ptr.Ref(50))
	assert.Equal(t, 50, *ptr.Ref(50))
	assert.Equal(t, "hello", *ptr.Ref("hello"))
}
