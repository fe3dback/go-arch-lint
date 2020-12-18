package view

// language=gohtml
const Mapping = `
	{{- /* gotype: github.com/fe3dback/go-arch-lint/internal/models.Mapping */ -}}
	{{ $root := .ProjectDirectory -}}

	module: {{ .ModuleName | colorize "green" }}
	Project Packages:
	{{ if eq .Scheme "list" -}}
		{{ $prev := "" -}}
		{{ range .MappingList -}}
			{{ $packageName := (.FileName | trimPrefix $root | dir | def "/") -}}

			{{ if ne $prev $packageName -}}
				{{ "  " }} {{ .ComponentName | padRight 20 " " -}}
				{{ $packageName | colorize "cyan" }}
			{{ end -}}
			
			{{ $prev = $packageName -}}
		{{ end -}}
	{{ else -}}
		{{ range .MappingGrouped -}}
			{{ "  " }} {{ .ComponentName }}:
			{{ $prev := "" -}}
			{{ range .FileNames -}}
				{{ $packageName := (. | trimPrefix $root | dir | def "/") -}}

				{{ if ne $prev $packageName -}}
					{{ "    " }} {{ $packageName | colorize "cyan" }}
				{{ end -}}

				{{ $prev = $packageName -}}
			{{ end -}}
		{{ end -}}
	{{ end -}}
`
