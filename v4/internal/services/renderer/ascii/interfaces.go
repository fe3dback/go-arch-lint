package ascii

type renderer interface {
	RegisterTemplate(id string, text []byte) error
	Render(id string, model any) (string, error)
}
