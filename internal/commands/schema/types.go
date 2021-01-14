package schema

type (
	JSONSchemaProvider interface {
		Provide(version int) (string, error)
	}
)
