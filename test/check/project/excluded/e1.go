package excluded

import (
	"github.com/fe3dback/go-arch-lint/test/check/project/a"
	"github.com/fe3dback/go-arch-lint/test/check/project/b"
)

func E1() {
	a.A1() // not allowed, but not checked by excluded dir
	b.B1() // not allowed, but not checked by excluded dir
}
