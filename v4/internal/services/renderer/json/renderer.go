package json

type JSON struct {
}

func NewRenderer() *JSON {
	return &JSON{}
}

func (r *JSON) Render(model any) (string, error) {
	return "todo", nil
}
