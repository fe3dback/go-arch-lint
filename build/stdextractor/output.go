package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type (
	outputDto struct {
		GeneratedAt   string
		PackageName   string
		GoVersionInfo string
		StdPackages   []outputStdPackage
	}

	outputStdPackage struct {
		Name string
	}
)

// language=gohtml
var outputTemplate = `// do not edit
// generated at {{.GeneratedAt}}
// generated for '{{.GoVersionInfo}}'

package {{.PackageName}}

var StdPackages = map[string]struct{}{
	{{range .StdPackages}}
	"{{.Name}}": {},
	{{- end}}
}
`

func renderOutput(dto outputDto) ([]byte, error) {
	tpl, err := template.
		New("template").
		Parse(outputTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed parse output template: %w", err)
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.Bytes(), nil
}
