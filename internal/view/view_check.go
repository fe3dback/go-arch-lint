package view

// language=gohtml
const Check = `
	module: {{.ModuleName | colorize "green"}}
	
	{{ range .DocumentNotices -}}
		{{.Text}}
		{{ if .SourceCodePreview -}}
			{{ .SourceCodePreview | printf "%s" -}}
		{{ end -}}
	{{ else -}}
		{{ if .ArchHasWarnings -}}
			{{ $warnCount := (plus (plus (len .ArchWarningsDependency) (len .ArchWarningsMatch)) (len .ArchWarningsDeepScan) ) -}}
			{{ range .ArchWarningsDependency -}}
				Component {{.ComponentName | colorize "magenta"}} file {{ .FileRelativePath | colorize "cyan"}} shouldn't depend on {{ .ResolvedImportName | colorize "blue"}}
				{{ if .SourceCodePreview -}}
					{{ .SourceCodePreview | printf "%s" -}}
				{{ end -}}
			{{ end -}}
			{{ range .ArchWarningsMatch -}}
				File {{.FileRelativePath | colorize "cyan"}} not attached to any component in archfile
				{{ if .SourceCodePreview -}}
					{{ .SourceCodePreview | printf "%s" -}}
				{{ end -}}
			{{ end }}
			{{ range .ArchWarningsDeepScan }}
				Dependency {{.Dependency.ComponentName | colorize "magenta"}} -\-> {{.Gate.ComponentName | colorize "magenta"}} not allowed
				  ├─ {{.Dependency.ComponentName | colorize "magenta"}} {{.Dependency.Name | colorize "blue"}} in {{ .Target.RelativePath | colorize "gray" }}
				  └─ {{.Gate.ComponentName | colorize "magenta"}} {{.Gate.MethodName | colorize "blue"}} in {{ .Gate.RelativePath | colorize "gray" }}
				{{ " " }}
				{{ concat "     " .Dependency.Injection.File ":" .Dependency.Injection.Line | colorize "gray" }}
				{{ if .Dependency.SourceCodePreview -}}
					{{ .Dependency.SourceCodePreview | printf "%s" | linePrefix "     " | colorize "red" -}}
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
