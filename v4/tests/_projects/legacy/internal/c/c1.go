package c

import "github.com/fe3dback/go-arch-lint/v4/tests/_projects/legacy/internal/a"

func C1() {
	a.A1() // not allowed
}
