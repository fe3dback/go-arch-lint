module github.com/fe3dback/go-arch-lint/test/check/project

go 1.18

require (
	github.com/example/a v0.0.0
	github.com/example/b v0.0.0
)

replace (
	github.com/example/a => ../stubs/example/a
	github.com/example/b => ../stubs/example/b
)
