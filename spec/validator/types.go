package validator

type (
	AnnotatedWarningParser interface {
		Parse(sourceText string) (line, pos int, err error)
	}

	PathResolver interface {
		Resolve(absPath string) (resolvePaths []string, err error)
	}
)
