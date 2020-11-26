package yamlresolver

type (
	AnnotationParser interface {
		Parse(annotation string) (line int, pos int, err error)
	}
)
