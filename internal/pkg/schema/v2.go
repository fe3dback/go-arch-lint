package schema

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
			"description": "Linter will prepend all path's in project with this relative path prefix (relative directory for analyse)",
			"type": "string"
		},
		"settings": {
			"title": "Global Scheme options",
			"type": "object",
			"additionalProperties": false,
			"properties": {
				"depOnAnyVendor": {
					"title": "allow import any vendor code to any project file",
					"type": "boolean"
				}
			}
		},
		"exclude": {
			"title": "Excluded folders from analyse",
			"type": "array",
			"items": {
				"type": "string",
				"title": "list of directories (relative path) for exclude from analyse"
			}
		},
		"excludeFiles": {
			"title": "Excluded files from analyse matched by regexp",
			"description": "package will by excluded in all package files is matched by provided regexp's",
			"type": "array",
			"items": {
				"type": "string",
				"title": "regular expression rules for file names, will exclude this files and it's packages from analyse",
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
					"anyOf": [
						{"$ref": "#/definitions/vendorIn"},
						{"type": "array", "items": {"$ref": "#/definitions/vendorIn"}}
					]
				}
			},
			"additionalProperties": false
		},
		"vendorIn": {
			"title": "full import path to vendor",
			"description": "one or more import path of vendor libs, support glob masking (src/\\*/engine/\\*\\*)",
			"type": "string",
			"examples": ["golang.org/x/mod/modfile", "example.com/*/libs/**", ["gopkg.in/yaml.v2", "github.com/mailru/easyjson"]]
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
					"anyOf": [
						{"$ref": "#/definitions/componentIn"},
						{"type": "array", "items": {"$ref": "#/definitions/componentIn"}}
					]
				}
			},
			"additionalProperties": false
		},
		"componentIn": {
			"title": "relative path to project package",
			"description": "relative directory name, support glob masking (src/\\*/engine/\\*\\*)",
			"type": "string",
			"examples": ["src/services", "src/services/*/repo", "src/*/services/**"]
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
					"description": "all component code can import any other project code, useful for DI/main component",
					"type": "boolean"
				},
				"anyVendorDeps": {
					"title": "Allow import any vendor package?",
					"description": "all component code can import any vendor code",
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
