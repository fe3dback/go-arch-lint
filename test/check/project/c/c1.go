package c

import "github.com/fe3dback/go-arch-lint/test/check/project/a"

func C1() {
	a.A1() // not allowed
}
