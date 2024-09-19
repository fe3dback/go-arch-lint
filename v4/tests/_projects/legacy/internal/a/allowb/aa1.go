package allowb

import (
	"github.com/fe3dback/go-arch-lint/v4/tests/_projects/legacy/internal/b"
	"github.com/fe3dback/go-arch-lint/v4/tests/_projects/legacy/internal/common/sub/foo/bar"
)

func AA1() {
	bar.BR1() // allowed common
	b.B1()    // allowed by deps
}
