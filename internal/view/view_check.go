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
			{{ range .ArchWarningsDependency -}}
				[WARN] Component '{{.ComponentName | colorize "green"}}': file '
				{{- .FileRelativePath | colorize "cyan"}}' shouldn't depend on '
				{{- .ResolvedImportName | colorize "yellow"}}'
				{{ if .SourceCodePreview -}}
					{{ .SourceCodePreview | printf "%s" -}}
				{{ end -}}
			{{ end -}}
			{{ range .ArchWarningsMatch -}}
				[WARN] File '{{.FileRelativePath | colorize "cyan"}}' not attached to any component in archfile
				{{ if .SourceCodePreview -}}
					{{ .SourceCodePreview | printf "%s" -}}
				{{ end -}}
			{{ end }}

			warnings found: {{len .ArchWarningsDependency | printf "%d" | colorize "yellow"}}
		{{ else -}}
			{{"OK - No warnings found" | colorize "green" -}}
		{{ end -}}
	{{ end -}}

	{{- /*gotype: github.com/fe3dback/go-arch-lint/internal/models.Check*/ -}}
`
