package scheme

// language=JSON
const V1 = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"id": "https://github.com/fe3dback/go-arch-lint/v1",
	"title": "Go Arch Lint V1",
	"type": "object",
	"description": "Arch file scheme version 1",
	"required": ["version", "components", "deps"],
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
		},
		"exclude": {
			"title": "Excluded folders from analyse",
			"type": "array",
			"items": {
				"type": "string",
				"title": "relative path to project root"
			}
		},
		"excludeFiles": {
			"title": "Excluded files from analyse matched by regexp",
			"description": "package will by excluded in all package files is matched by provided regexp's",
			"type": "array",
			"items": {
				"type": "string",
				"title": "regular expression for absolute file path matching",
				"x-intellij-language-injection": "regexp"
			}
		},
		"vendors": {
			"title": "List of vendors libs",
			"type": "object",
			"additionalProperties": {
				"type": "object",
				"required": ["in"],
				"properties": {
					"in": {
						"title": "full import path to vendor",
						"type": "string",
						"examples": ["golang.org/x/mod/modfile"]
					}
				},
				"additionalProperties": false
			}
		},
		"components": {
			"title": "List of components",
			"type": "object",
			"additionalProperties": {
				"type": "object",
				"required": ["in"],
				"properties": {
					"in": {
						"title": "relative path to project package",
						"description": "can contain glob for search",
						"type": "string",
						"examples": ["src/services", "src/services/*/repo", "src/*/services/**"]
					}
				},
				"additionalProperties": false
			}
		},
		"commonComponents": {
			"title": "List of components names",
			"description": "All project packages can import this components, useful for utils packages like 'models'",
			"type": "array",
			"items": {
				"type": "string",
				"title": "component name"
			}
		},
		"commonVendors": {
			"title": "List of vendor names",
			"description": "All project packages can import this vendor libs",
			"type": "array",
			"items": {
				"type": "string",
				"title": "vendor name"
			}
		},
		"deps": {
			"title": "",
			"type": "object",
			"additionalProperties": {
				"type": "object",
				"properties": {
					"anyProjectDeps": {
						"title": "Allow import any project package?",
						"type": "boolean"
					},
					"anyVendorDeps": {
						"title": "Allow import any vendor package?",
						"type": "boolean"
					},
					"mayDependOn": {
						"title": "List of allowed components to import",
						"type": "array",
						"items": {
							"type": "string",
							"title": "component name"
						}
					},
					"canUse": {
						"title": "List of allowed vendors to import",
						"type": "array",
						"items": {
							"type": "string",
							"title": "vendor name"
						}
					}
				},
				"additionalProperties": false
			}
		}
	},
	"additionalProperties": false
}`
