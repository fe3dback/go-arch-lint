package schema

type (
	jsonSchemaProvider interface {
		Provide(version int) (string, error)
	}
)
