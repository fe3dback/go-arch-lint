package models

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
