package view

// language=gohtml
const Version = `
	Linter version: {{.LinterVersion | colorize "yellow" }}
	Supported go arch file versions: {{.GoArchFileSupported | colorize "yellow" }}
	Build time: {{.BuildTime | colorize "yellow" }}
	Commit hash: {{.CommitHash | colorize "yellow" }}

	{{- /*gotype: github.com/fe3dback/go-arch-lint/internal/models.Version*/ -}}
`
