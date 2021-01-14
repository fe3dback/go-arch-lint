package b

import "github.com/fe3dback/go-arch-lint/test/check/project/internal/common/sub/foo/bar"

func B1() {
	bar.BR1() // common - allowed
}
