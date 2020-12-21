package scheme

// language=JSON
const V1 = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"id": "https://github.com/fe3dback/go-arch-lint/v1",
	"title": "Go Arch Lint V1",
	"type": "object",
	"description": "Arch file scheme version 1",
	"required": ["version"],
	"properties": {
		"version": {
			"title": "Scheme Version",
			"description": "Defines arch file syntax and file validation rules",
			"type": "integer",
			"minimum": 1,
			"maximum": 1
		},
		"allow": {
			"title": "Global Scheme options",
			"type": "object",
			"properties": {
				"depOnAnyVendor": {
					"title": "Any project file can import any vendor lib",
					"type": "boolean"
				}
			},
			"additionalProperties": false
		}
	},
	"additionalProperties": false
}`
