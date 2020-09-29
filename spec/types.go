package spec

type (
	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
