package foo

import (
	"fmt"
)

type Foo []string

func FooF(message string, args ...any) Foo {
	return Foo{fmt.Sprintf(message, args)}
}
