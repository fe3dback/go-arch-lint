package models

type (
	ArchFileSyntaxWarning struct {
		srcPath  string
		srcError error
	}
)

func NewArchFileSyntaxWarning(path string, err error) *ArchFileSyntaxWarning {
	return &ArchFileSyntaxWarning{
		srcPath:  path,
		srcError: err,
	}
}

func (w ArchFileSyntaxWarning) Path() string {
	return w.srcPath
}

func (w ArchFileSyntaxWarning) Warning() error {
	return w.srcError
}
