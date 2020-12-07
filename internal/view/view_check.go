package view

// language=gohtml
const Check = `
	module: {{.ModuleName | colorize "green"}}
	
	{{ range .DocumentNotices -}}
		[Archfile] {{.Text}}
		{{ if .SourceCodePreview -}}
			{{ .SourceCodePreview | printf "%s" -}}
		{{ end -}}
	{{ else -}}
		{{ if .ArchHasWarnings -}}
			// has warnings
		{{ else -}}
			{{"OK - No warnings found" | colorize "green" -}}
		{{ end -}}
	{{ end -}}

	{{- /*gotype: github.com/fe3dback/go-arch-lint/internal/models.Check*/ -}}
`
