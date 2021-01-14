package allowb

import (
	"github.com/fe3dback/go-arch-lint/test/check/project/internal/b"
	"github.com/fe3dback/go-arch-lint/test/check/project/internal/common/sub/foo/bar"
)

func AA1() {
	bar.BR1() // allowed common
	b.B1()    // allowed by deps
}
