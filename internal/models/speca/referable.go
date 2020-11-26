package speca

type (
	Referable interface {
		Reference() Reference
	}

	Reference struct {
		Valid  bool
		File   string
		Line   int
		Offset int
	}
)

func NewReference(file string, line int, offset int) Reference {
	return Reference{Valid: true, File: file, Line: line, Offset: offset}
}

func NewEmptyReference() Reference {
	return Reference{Valid: false, File: "", Line: 0, Offset: 0}
}
