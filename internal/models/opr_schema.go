package models

type (
	CmdSchemaIn struct {
		Version int
	}

	CmdSchemaOut struct {
		Version    int    `json:"Version"`
		JSONSchema string `json:"JsonSchema"`
	}
)
