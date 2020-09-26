package validator

type (
	Warning struct {
		yamlPath    string
		yamlWarning error
	}
)

func (w Warning) Path() string {
	return w.yamlPath
}

func (w Warning) Warning() error {
	return w.yamlWarning
}
