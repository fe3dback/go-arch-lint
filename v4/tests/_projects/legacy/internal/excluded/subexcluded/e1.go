package subexcluded

import (
	"github.com/fe3dback/go-arch-lint/v4/tests/_projects/legacy/internal/a"
	"github.com/fe3dback/go-arch-lint/v4/tests/_projects/legacy/internal/b"
)

func E1() {
	a.A1() // not allowed, but not checked by excluded dir
	b.B1() // not allowed, but not checked by excluded dir
}
