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
		Hint   string // can contain error text or other hint information in case of Valid=false
	}
)
