package build

import "time"

var (
	Version     = "dev"      // nolint
	CompileTime = time.Now() // nolint
	CommitHash  = "unknown"  // nolint
)
