package view

// language=gohtml
const Check = `
	module: {{.ModuleName | colorize "green"}}

	{{ range .DocumentNotices -}}
	[Archfile] {{.Text}}
	{{ if .SourceCodePreview -}}
		{{.SourceCodePreview | printf "%s"}}
	{{- end -}}
	{{else}}
		// not have notice
	{{end}}

	{{- /*gotype: github.com/fe3dback/go-arch-lint/internal/models.Check*/ -}}
`
