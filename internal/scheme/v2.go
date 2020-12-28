package scheme

// language=JSON
const V2 = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"id": "https://github.com/fe3dback/go-arch-lint/v2",
	"title": "Go Arch Lint V2",
	"type": "object",
	"description": "Arch file scheme version 2",
	"required": ["version", "components", "deps"],
	"additionalProperties": false,
	"properties": {
		"version": {"$ref": "#/definitions/version"},
		"workdir": {"$ref": "#/definitions/workdir"},
		"allow": {"$ref": "#/definitions/settings"},
		"exclude": {"$ref": "#/definitions/exclude"},
		"excludeFiles": {"$ref": "#/definitions/excludeFiles"},
		"vendors": {"$ref": "#/definitions/vendors"},
		"commonVendors": {"$ref": "#/definitions/commonVendors"},
		"components": {"$ref": "#/definitions/components"},
		"commonComponents": {"$ref": "#/definitions/commonComponents"},
		"deps": {"$ref": "#/definitions/dependencies"}
	},
	"definitions": {
		"version": {
			"title": "Scheme Version",
			"description": "Defines arch file syntax and file validation rules",
			"type": "integer",
			"minimum": 2,
			"maximum": 2
		},
		"workdir": {
			"title": "Working directory",
			"description": "Linter will prepend all path's in project with this relative path prefix",
			"type": "string"
		},
		"settings": {
			"title": "Global Scheme options",
			"type": "object",
			"additionalProperties": false,
			"properties": {
				"depOnAnyVendor": {
					"title": "Any project file can import any vendor lib",
					"type": "boolean"
				}
			}
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
			"title": "List of vendor libs",
			"type": "object",
			"additionalProperties": {"$ref": "#/definitions/vendor"}
		},
		"vendor": {
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
		"components": {
			"title": "List of components",
			"type": "object",
			"additionalProperties": {"$ref": "#/definitions/component"}
		},
		"component": {
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
		"dependencies": {
			"title": "Dependency rules between spec and package imports",
			"type": "object",
			"additionalProperties": {"$ref": "#/definitions/dependencyRule"}
		},
		"dependencyRule": {
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
}`
