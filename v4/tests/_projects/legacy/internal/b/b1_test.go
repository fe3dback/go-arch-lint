package b

import (
	"testing"

	"github.com/fe3dback/go-arch-lint/v4/tests/_projects/legacy/internal/a"
)

func TestB1(t *testing.T) {
	a.A1() // not allowed, but not checked (excluded by regexp)
}
