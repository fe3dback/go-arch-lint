package e

import (
	libA "github.com/example/a"
	libB "github.com/example/b"

	modelA "github.com/fe3dback/go-arch-lint/test/check/project/internal/d/models/a/model"
	modelB "github.com/fe3dback/go-arch-lint/test/check/project/internal/d/models/b/model"
)

func E1() {
	modelA.ModelA()
	modelB.ModelB()

	libA.LibraryA()
	libB.LibraryB()
}
