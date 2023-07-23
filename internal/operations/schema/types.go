package schema

type (
	jsonSchemaProvider interface {
		Provide(version int) ([]byte, error)
	}
)
