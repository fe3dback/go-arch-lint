package b

import (
	"testing"

	"github.com/fe3dback/go-arch-lint/test/check/project/a"
)

func TestB1(t *testing.T) {
	a.A1() // not allowed, but not checked (excluded by regexp)
}
