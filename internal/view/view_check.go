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
			{{ $warnCount := (plus (plus (len .ArchWarningsDependency) (len .ArchWarningsMatch)) (len .ArchWarningsDeepScan) ) -}}
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
			{{ range .ArchWarningsDeepScan }}
				Dependency {{.Dependency.ComponentName | colorize "green"}} -/-> {{.Gate.ComponentName | colorize "green"}} not allowed
				  ├─ {{.Dependency.ComponentName | colorize "green"}}: type {{.Dependency.Name | colorize "cyan"}} injected in {{.Dependency.InjectionPath | colorize "gray"}}
				  └─ {{.Gate.ComponentName | colorize "green"}}: into {{.Gate.MethodName | colorize "cyan"}} in {{.Gate.RelativePath | colorize "gray"}}
				  
				{{ if .Dependency.SourceCodePreview -}}
					{{ .Dependency.SourceCodePreview | printf "%s" -}}
				{{ end -}}
				  
			{{ end }}

			--
			total notices: {{ plus $warnCount .OmittedCount | printf "%d" | colorize "yellow" }}
			{{ if gt .OmittedCount 0 -}}
				omitted: {{.OmittedCount | printf "%d" | colorize "yellow" }} (too big to display)
			{{ end }}
		{{ else -}}
			{{"OK - No warnings found" | colorize "green" -}}
		{{ end -}}
	{{ end -}}

	{{- /*gotype: github.com/fe3dback/go-arch-lint/internal/models.Check*/ -}}
`
