package schema

type (
	JsonSchemaProvider interface {
		Provide(version int) (string, error)
	}
)
