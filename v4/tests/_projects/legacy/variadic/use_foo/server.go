package use_foo

import (
	"github.com/fe3dback/go-arch-lint/v4/tests/_projects/legacy/variadic/foo"
)

var (
	foo1 = foo.FooF("boo")
	foo2 = foo.FooF("boo", 5)
	foo3 = foo.FooF("boo", 5, "test", 15)
)
